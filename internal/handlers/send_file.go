package handlers

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"net/textproto"
	"strings"
)

var allowedMimeTypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/pdf": true,
}

func (h *Handler) HandleSendFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(h.Config.BufUploadSizeMail)
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["file"]
	if files == nil {
		http.Error(w, "No files provided", http.StatusBadRequest)
		return
	}

	if len(files) > 1 {
		http.Error(w, "Please attach only one file", http.StatusBadRequest)
		return
	}

	fileHeader := files[0]

	// file, fileHeader, err := r.FormFile("file")
	// if err != nil {
	// 	http.Error(w, "Failed to get file", http.StatusBadRequest)
	// 	return
	// }
	// defer file.Close()

	if !allowedMimeTypes[fileHeader.Header.Get("Content-Type")] {
		http.Error(w, "Unsupported file type", http.StatusBadRequest)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Получение списка email-адресов из формы
	emailList := strings.Split(r.FormValue("emails"), ",")
	if len(emailList) == 0 || emailList[0] == "" {
		http.Error(w, "No valid emails provided", http.StatusBadRequest)
		return
	}

	// Чтение первых 512 байт для определения MIME типа
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	expectedType := fileHeader.Header.Get("Content-Type")
	detectedType := http.DetectContentType(buf)
	if detectedType != expectedType {
		http.Error(w, "Actual file type does not match expected file type", http.StatusBadRequest)
		return
	}

	// Возвращение указателя в начало файла перед отправкой
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, "Failed to seek file", http.StatusInternalServerError)
		return
	}

	if err := h.sendEmailWithAttachment(file, fileHeader.Filename, emailList, fileHeader.Header.Get("Content-Type")); err != nil {
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) sendEmailWithAttachment(file multipart.File, filename string, recipients []string, contentType string) error {
	// Настройка информации для аутентификации
	auth := smtp.PlainAuth("", h.Config.SMTPUser, h.Config.SMTPPassword, h.Config.SMTPHost)

	// Создание буфера для письма
	var email bytes.Buffer
	writer := multipart.NewWriter(&email)
	defer writer.Close()

	// Добавление заголовка
	headers := make(map[string]string)
	headers["From"] = "github.com/asari92"
	headers["To"] = strings.Join(recipients, ",")
	headers["Subject"] = "Document"
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "multipart/mixed; boundary=" + writer.Boundary()
	writeHeaders(&email, headers)

	// Добавление тела письма
	part, err := writer.CreatePart(textproto.MIMEHeader{"Content-Type": {"text/plain; charset=utf-8"}})
	if err != nil {
		return err
	}
	_, err = part.Write([]byte("Here is the document you requested.\n"))
	if err != nil {
		return err
	}

	// Добавление файла как вложения
	part, err = writer.CreatePart(textproto.MIMEHeader{
		"Content-Disposition":       {"attachment; filename=\"" + filename + "\""},
		"Content-Type":              {contentType},
		"Content-Transfer-Encoding": {"base64"},
	})
	if err != nil {
		return err
	}
	encoder := base64.NewEncoder(base64.StdEncoding, part)
	if _, err = io.Copy(encoder, file); err != nil {
		return err
	}
	encoder.Close()

	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!TODO  переделать с использование горутин и waitgroup чтобы каждому адресату отдельно обрабатывалось отправка письма
	// Отправка письма
	return smtp.SendMail(h.Config.SMTPHost+":"+h.Config.SMTPPort, auth, "github.com/asari92", recipients, email.Bytes())
}

func writeHeaders(w io.Writer, headers map[string]string) {
	for k, v := range headers {
		w.Write([]byte(k + ": " + v + "\r\n"))
	}
	w.Write([]byte("\r\n"))
}
