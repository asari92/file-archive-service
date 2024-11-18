package handlers

import (
	"encoding/json"
	"net/http"

	"file-archive-service/pkg/validator"
)

var archiverAllowedMimeTypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/xml": true,
	"image/jpeg":      true,
	"image/png":       true,
}

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

	err = validator.ValidateMimeTypes(files, archiverAllowedMimeTypes)
	if err != nil {
		app.Logger.Error("ValidateMimeTypes", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
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
