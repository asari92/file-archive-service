package handlers

import (
	"archive/zip"
	"encoding/json"
	"net/http"

	"file-archive-service/pkg/utils"
)

func (app *Application) HandleArchiveInformation(w http.ResponseWriter, r *http.Request) {
	// Парсинг формы с файлом
	err := r.ParseMultipartForm(app.Config.BufUploadSizeInfo)
	if err != nil {
		http.Error(w, "Error parsing multipart/form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	mFile, mFileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get the file", http.StatusBadRequest)
		return
	}
	defer mFile.Close()

	// Проверяем, что файл является zip-архивом
	zipReader, err := zip.NewReader(mFile, mFileHeader.Size)
	if err != nil {
		http.Error(w, "File is not a valid zip archive", http.StatusBadRequest)
		return
	}

	// Сбор информации о файлах в архиве
	var files []FileInfo
	var totalSize float64
	for _, f := range zipReader.File {
		info := f.FileHeader.FileInfo()

		// Используем функцию getMimeType для получения MIME-типа файла
		mimeType := utils.GetMimeType(f.Name)

		file := FileInfo{
			FilePath: f.Name,
			Size:     float64(info.Size()),
			MimeType: mimeType,
		}
		files = append(files, file)
		totalSize += float64(info.Size())
	}

	// Формирование и отправка ответа
	response := Response{
		Filename:    mFileHeader.Filename,
		ArchiveSize: float64(mFileHeader.Size),
		TotalSize:   totalSize,
		TotalFiles:  len(files),
		Files:       files,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
