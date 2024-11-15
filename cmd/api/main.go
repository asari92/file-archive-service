package main

import (
	"log"
	"net/http"
	"os"

	"file-archive-service/internal/handlers"
	"file-archive-service/internal/utils"
)

func main() {
	// Загрузите переменные окружения из файла .env
	if err := utils.LoadEnv(".env"); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	// Используйте переменные окружения в вашем приложении
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080" // Значение по умолчанию
	}

	mux := http.NewServeMux()
	mux.Handle("POST /api/archive/information", http.HandlerFunc(handlers.HandleArchiveInformation))

	log.Printf("Server starting on port %s", port)
	// Запуск сервера
	log.Fatal(http.ListenAndServe(port, mux)) // `port` уже содержит нужное двоеточие
}
