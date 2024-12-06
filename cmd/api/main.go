package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"file-archive-service/internal/domain/archive"
	"file-archive-service/internal/domain/mail"
	"file-archive-service/internal/handler"
	"file-archive-service/internal/service"
	"file-archive-service/pkg/config"
	"file-archive-service/pkg/utils"
)

func main() {
	// Загрузите переменные окружения из файла .env
	if err := utils.LoadEnv(".env"); err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}

	conf := config.New()

	// Инициализация нового логгера slog
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,            // Включить вывод источника вызова (файл и строка)
		Level:     slog.LevelDebug, // задан дебаг уровень, можно поменять на инфо чтобы убрать лишнюю инфу
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				// Устанавливаем формат времени на "2006-01-02 15:04:05"
				a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05"))
			}
			return a
		},
	}))

	slog.SetDefault(logger)

	mailer := mail.NewGoMailAdapter(conf)

	h := &handler.Handler{
		Config:  conf,
		Logger:  logger,
		Service: service.NewService(archive.NewZipArchiver(), mailer, conf),
	}

	err := serve(h, conf.Host+":"+conf.Port)
	if err != nil {
		logger.Error("Fatal server error", "error", err)
	}
}

func serve(h *handler.Handler, addr string) error {
	srv := &http.Server{
		Addr:         addr,
		ErrorLog:     slog.NewLogLogger(h.Logger.Handler(), slog.LevelError),
		Handler:      h.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  9 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Create a shutdownError channel. We will use this to receive any errors returned
	// by the graceful Shutdown() function.
	shutdownError := make(chan error)

	go func() {
		// Intercept the signals, as before.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		// Update the log entry to say "shutting down server" instead of "caught signal".
		h.Logger.Info("shutting down server", "signal", s.String())

		// Create a context with a 20-second timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// Call Shutdown() on our server, passing in the context we just made.
		// Shutdown() will return nil if the graceful shutdown was successful, or an
		// error (which may happen because of a problem closing the listeners, or
		// because the shutdown didn't complete before the 20-second context deadline is
		// hit). We relay this return value to the shutdownError channel.
		shutdownError <- srv.Shutdown(ctx)
	}()

	h.Logger.Info("Starting server", "address", addr)
	// Calling Shutdown() on our server will cause ListenAndServe() to immediately
	// return a http.ErrServerClosed error. So if we see this error, it is actually a
	// good thing and an indication that the graceful shutdown has started. So we check
	// specifically for this, only returning the error if it is NOT http.ErrServerClosed.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Otherwise, we wait to receive the return value from Shutdown() on the
	// shutdownError channel. If return value is an error, we know that there was a
	// problem with the graceful shutdown and we return the error.
	err = <-shutdownError
	if err != nil {
		return err
	}

	// At this point we know that the graceful shutdown completed successfully and we
	// log a "stopped server" message.
	h.Logger.Info("stopped server", "addr", addr)

	return nil
}
