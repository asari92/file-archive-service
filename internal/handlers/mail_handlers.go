package handlers

import (
	"net/http"
	"strings"

	"file-archive-service/pkg/validator"
)

var mailerAllowedSignatures = map[string][]byte{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {0x50, 0x4B, 0x03, 0x04}, // DOCX
	"application/pdf": {0x25, 0x50, 0x44, 0x46}, // PDF
}

func (app *Application) HandleSendFile(w http.ResponseWriter, r *http.Request) {
	// Парсинг формы и базовая проверка
	if err := r.ParseMultipartForm(app.Config.BufUploadSizeMail); err != nil {
		app.Logger.Error("Error parsing form", "error", err)
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["file"]
	if files == nil || len(files) != 1 {
		app.Logger.Error("zero or more than 1 file")
		http.Error(w, "Please attach only one file", http.StatusBadRequest)
		return
	}

	fileHeader := files[0]

	if fileHeader.Size >= app.Config.MaxSendFileSize {
		app.Logger.Error("file size exceeds the maximum limit")
		http.Error(w, "file size exceeds the maximum limit", http.StatusBadRequest)
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		app.Logger.Error("failed to open file", "error", err)
		http.Error(w, "failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if err := validator.ValidateFileSignature(file, fileHeader.Header.Get("Content-Type"), mailerAllowedSignatures); err != nil {
		app.Logger.Error("ValidateFileSignature", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получение и фильтрация email-адресов
	emailList := strings.Split(r.FormValue("emails"), ",")
	recipients, err := validator.ValidateEmails(emailList)
	if err != nil {
		app.Logger.Error("ValidateEmails", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Подготовка и отправка писем
	err = app.Service.Mailer.SendEmailWithAttachment(app.Config.MailFrom, recipients, "Document", fileHeader.Filename, "", file)
	if err != nil {
		app.Logger.Error("SendEmailWithAttachment", "error", err)
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
	}

	// поскольку по ТЗ должен возвращаться http.StatusOK, значит нужно дождаться пока письма отправятся.
	// если отправить письма в горутине и не дожидаться ответа нужно возращать http.StatusAccepted
	w.WriteHeader(http.StatusOK)
}
