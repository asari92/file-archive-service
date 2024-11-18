package handlers

import (
	"encoding/json"
	"net/http"
)

// func (app *Application) HandleCreateArchive(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseMultipartForm(app.Config.BufUploadSizeCreate)
// 	if err != nil {
// 		http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	files := r.MultipartForm.File["files[]"]
// 	if files == nil {
// 		http.Error(w, "No files provided", http.StatusBadRequest)
// 		return
// 	}

// 	// Создание архива
// 	buf := new(bytes.Buffer)
// 	zipWriter := zip.NewWriter(buf)

// 	// Обработка каждого файла
// 	for _, fileHeader := range files {
// 		// Проверяем MIME-тип
// 		if !allowedMIMETypes[fileHeader.Header.Get("Content-Type")] {
// 			http.Error(w, "Invalid file type: "+fileHeader.Filename, http.StatusBadRequest)
// 			return
// 		}

// 		file, err := fileHeader.Open()
// 		if err != nil {
// 			http.Error(w, "Error opening file: "+fileHeader.Filename, http.StatusInternalServerError)
// 			return
// 		}

// 		// Добавляем файл в архив
// 		zipFile, err := zipWriter.Create(fileHeader.Filename)
// 		if err != nil {
// 			http.Error(w, "Error creating zip entry: "+fileHeader.Filename, http.StatusInternalServerError)
// 			return
// 		}

// 		_, err = io.Copy(zipFile, file)
// 		file.Close() // Закрываем файл после копирования его содержимого
// 		if err != nil {
// 			http.Error(w, "Error writing to zip: "+fileHeader.Filename, http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	// Закрытие архива
// 	err = zipWriter.Close()
// 	if err != nil {
// 		http.Error(w, "Error closing zip writer", http.StatusInternalServerError)
// 		return
// 	}

// 	// Устанавливаем заголовки ответа
// 	w.Header().Set("Content-Type", "application/zip")
// 	// w.Header().Set("Content-Disposition", "attachment; filename=\"archive.zip\"")

// 	// Отправка архива
// 	w.Write(buf.Bytes())
// }

// func (app *Application) HandleArchiveInformation(w http.ResponseWriter, r *http.Request) {
// 	// Парсинг формы с файлом
// 	err := r.ParseMultipartForm(app.Config.BufUploadSizeInfo)
// 	if err != nil {
// 		http.Error(w, "Error parsing multipart/form data: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	mFile, mFileHeader, err := r.FormFile("file")
// 	if err != nil {
// 		http.Error(w, "Failed to get the file", http.StatusBadRequest)
// 		return
// 	}
// 	defer mFile.Close()

// 	// Проверяем, что файл является zip-архивом
// 	zipReader, err := zip.NewReader(mFile, mFileHeader.Size)
// 	if err != nil {
// 		http.Error(w, "File is not a valid zip archive", http.StatusBadRequest)
// 		return
// 	}

// 	// Сбор информации о файлах в архиве
// 	var files []FileInfo
// 	var totalSize float64
// 	for _, f := range zipReader.File {
// 		info := f.FileHeader.FileInfo()

// 		// Используем функцию getMimeType для получения MIME-типа файла
// 		mimeType := utils.GetMimeType(f.Name)

// 		file := FileInfo{
// 			FilePath: f.Name,
// 			Size:     float64(info.Size()),
// 			MimeType: mimeType,
// 		}
// 		files = append(files, file)
// 		totalSize += float64(info.Size())
// 	}

// 	// Формирование и отправка ответа
// 	response := Response{
// 		Filename:    mFileHeader.Filename,
// 		ArchiveSize: float64(mFileHeader.Size),
// 		TotalSize:   totalSize,
// 		TotalFiles:  len(files),
// 		Files:       files,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

func (app *Application) HandleCreateArchive(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(app.Config.BufUploadSizeCreate)
	if err != nil {
		app.Logger.Error("parsing multipart form", "error", err)
		http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files[]"]
	if files == nil {
		app.Logger.Error("no files provided", "error", err)
		http.Error(w, "No files provided", http.StatusBadRequest)
		return
	}

	buf, err := app.Service.CreateArchive(files)
	if err != nil {
		app.Logger.Error("CreateArchive", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Write(buf.Bytes())
}

func (app *Application) HandleArchiveInformation(w http.ResponseWriter, r *http.Request) {
	// 	// Парсинг формы с файлом
	err := r.ParseMultipartForm(app.Config.BufUploadSizeInfo)
	if err != nil {
		app.Logger.Error("parsing multipart/form data", "error", err)
		http.Error(w, "Error parsing multipart/form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	mFile, mFileHeader, err := r.FormFile("file")
	if err != nil {
		app.Logger.Error("failed to get the file", "error", err)
		http.Error(w, "Failed to get the file", http.StatusBadRequest)
		return
	}
	defer mFile.Close()

	response, err := app.Service.GenerateArchiveInfo(&mFile, mFileHeader)
	if err != nil {
		app.Logger.Error("GenerateArchiveInfo", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		app.Logger.Error("json encode", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
}
