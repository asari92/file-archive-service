package handlers

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
)

// Разрешенные MIME-типы файлов для архивации
var allowedMIMETypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/xml": true,
	"image/jpeg":      true,
	"image/png":       true,
}

func (app *Application) HandleCreateArchive(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(app.Config.BufUploadSizeCreate)
	if err != nil {
		http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files[]"]
	if files == nil {
		http.Error(w, "No files provided", http.StatusBadRequest)
		return
	}

	// Создание архива
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Обработка каждого файла
	for _, fileHeader := range files {
		// Проверяем MIME-тип
		if !allowedMIMETypes[fileHeader.Header.Get("Content-Type")] {
			http.Error(w, "Invalid file type: "+fileHeader.Filename, http.StatusBadRequest)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Error opening file: "+fileHeader.Filename, http.StatusInternalServerError)
			return
		}

		// Добавляем файл в архив
		zipFile, err := zipWriter.Create(fileHeader.Filename)
		if err != nil {
			http.Error(w, "Error creating zip entry: "+fileHeader.Filename, http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(zipFile, file)
		file.Close() // Закрываем файл после копирования его содержимого
		if err != nil {
			http.Error(w, "Error writing to zip: "+fileHeader.Filename, http.StatusInternalServerError)
			return
		}
	}

	// Закрытие архива
	err = zipWriter.Close()
	if err != nil {
		http.Error(w, "Error closing zip writer", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки ответа
	w.Header().Set("Content-Type", "application/zip")
	// w.Header().Set("Content-Disposition", "attachment; filename=\"archive.zip\"")

	// Отправка архива
	w.Write(buf.Bytes())
}
