package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"
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

	file, err := fileHeader.Open()
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if err := checkFileType(file, fileHeader.Header.Get("Content-Type")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получение списка email-адресов из формы
	emailList := strings.Split(r.FormValue("emails"), ",")
	if len(emailList) == 0 || emailList[0] == "" {
		http.Error(w, "No valid emails provided", http.StatusBadRequest)
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

	// Добавление файла как вложения
	part, err := writer.CreatePart(textproto.MIMEHeader{
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
	headers["Date"] = time.Now().Format(time.RFC1123Z)
	for k, v := range headers {
		w.Write([]byte(k + ": " + v + "\r\n"))
	}
	w.Write([]byte("\r\n"))
}

func checkFileType(file multipart.File, expectedType string) error {
	if !allowedMimeTypes[expectedType] {
		return fmt.Errorf("unsupported file type")
	}

	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	detectedType := http.DetectContentType(buf)
	if detectedType != expectedType {
		return fmt.Errorf("actual file type '%s' does not match expected file type '%s'", detectedType, expectedType)
	}

	// Возвращение указателя в начало файла
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek file: %v", err)
	}
	return nil
}
