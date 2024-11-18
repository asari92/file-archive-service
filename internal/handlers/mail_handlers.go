package handlers

import (
	"net/http"
	"strings"
)

func (app *Application) HandleSendFile(w http.ResponseWriter, r *http.Request) {
	// Парсинг формы и базовая проверка
	if err := r.ParseMultipartForm(app.Config.BufUploadSizeMail); err != nil {
		app.Logger.Error("Error parsing form", "error", err)
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["file"]
	if files == nil || len(files) != 1 {
		http.Error(w, "Please attach only one file", http.StatusBadRequest)
		return
	}

	fileHeader := files[0]
	emailList := strings.Split(r.FormValue("emails"), ",")

	// Передача логики в сервис
	if err := app.Service.ProcessAndSendFile(fileHeader, emailList); err != nil {
		app.Logger.Error("Error send file", "error", err)
		http.Error(w, "Failed to process request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// поскольку по ТЗ должен возвращаться http.StatusOK, значит нужно дождаться пока письма отправятся.
	// если отправить письма в горутине и не дожидаться ответа нужно возращать http.StatusAccepted
	w.WriteHeader(http.StatusOK)
}
