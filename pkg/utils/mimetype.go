package utils

import (
	"path/filepath"
	"strings"
)

// getMimeType возвращает MIME-тип на основе расширения файла
func GetMimeType(fileName string) string {
	// Получаем расширение файла
	ext := strings.ToLower(filepath.Ext(fileName))

	// Словарь расширений и соответствующих MIME-типов
	mimeTypes := map[string]string{
		".jpeg": "image/jpeg",
		".jpg":  "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".bmp":  "image/bmp",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".pdf":  "application/pdf",
		".zip":  "application/zip",
		".txt":  "text/plain",
		".xml":  "application/xml",
		// Добавьте дополнительные расширения и MIME-типы по мере необходимости
	}

	// Возвращаем MIME-тип, если он найден; иначе возвращаем стандартный тип
	if mimeType, ok := mimeTypes[ext]; ok {
		return mimeType
	}
	return "application/octet-stream"
}
