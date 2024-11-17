package main

import (
	"log"
	"log/slog"
	"os"

	"file-archive-service/internal/handlers"
	"file-archive-service/internal/service"
	"file-archive-service/internal/utils"
	"file-archive-service/pkg/config"
)

func main() {
	// Загрузите переменные окружения из файла .env
	if err := utils.LoadEnv(".env"); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
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

	app := &handlers.Application{
		Config:  conf,
		Logger:  logger,
		Service: service.NewService(conf),
	}

	err := app.Serve(conf.Host + ":" + conf.Port)
	if err != nil {
		logger.Error("Fatal server error", "error", err)
	}
}
