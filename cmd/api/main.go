package main

import (
	"log"
	"net/http"

	"file-archive-service/internal/handlers"
	"file-archive-service/internal/utils"
	"file-archive-service/pkg/config"
)

func main() {
	// Загрузите переменные окружения из файла .env
	if err := utils.LoadEnv(".env"); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	conf := config.New()

	handler := handlers.NewHandler(conf)

	mux := http.NewServeMux()
	mux.Handle("POST /api/archive/information", http.HandlerFunc(handler.HandleArchiveInformation))
	mux.Handle("POST /api/archive/files", http.HandlerFunc(handler.HandleCreateArchive))

	log.Printf("Server starting on port %s", conf.Port)
	// Запуск сервера
	log.Fatal(http.ListenAndServe(conf.Port, mux)) // `port` уже содержит нужное двоеточие
}
