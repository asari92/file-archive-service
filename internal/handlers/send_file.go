package handlers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"

	gomail "gopkg.in/mail.v2"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

var allowedSignatures = map[string][]byte{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {0x50, 0x4B, 0x03, 0x04}, // DOCX
	"application/pdf": {0x25, 0x50, 0x44, 0x46}, // PDF
}

func (h *Handler) HandleSendFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(h.Config.BufUploadSizeMail)
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["file"]
	if files == nil || len(files) > 1 {
		http.Error(w, "Please attach only one file", http.StatusBadRequest)
		return
	}

	fileHeader := files[0]

	if fileHeader.Size >= h.Config.MaxSendFileSize {
		http.Error(w, "file size exceeds the maximum limit", http.StatusBadRequest)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		http.Error(w, "failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if err := checkFileSignature(file, fileHeader.Header.Get("Content-Type")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получение и фильтрация email-адресов
	emailList := strings.Split(r.FormValue("emails"), ",")
	recipients := validateEmails(emailList)
	if len(recipients) == 0 {
		http.Error(w, "No valid emails provided", http.StatusBadRequest)
		return
	}

	// Подготовка и отправка писем
	if err := h.sendEmailWithAttachment(file, fileHeader.Filename, recipients); err != nil {
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// поскольку по ТЗ должен возвращаться http.StatusOK, значит нужно дождаться пока письма отправятся.
	// если отправить письма в горутине и не дожидаться ответа нужно возращать http.StatusAccepted
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) sendEmailWithAttachment(file multipart.File, filename string, recipients []string) error {
	d := gomail.NewDialer(h.Config.SMTPHost, h.Config.SMTPPort, h.Config.SMTPUser, h.Config.SMTPPassword)
	d.Timeout = h.Config.DialerTimeout * time.Second // Устанавливаем таймаут на 60 секунд

	m := gomail.NewMessage()
	// m.SetHeader("From", h.Config.SMTPUser)
	m.SetHeader("From", "dulmaev.andrei@gmail.com")
	m.SetHeader("To", recipients...)
	m.SetHeader("Subject", "Document")
	m.AttachReader(filename, file)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Failed to send email:", err)
		return err
	}

	return nil
}

func checkFileSignature(file multipart.File, expectedType string) error {
	expectedSignature, ok := allowedSignatures[expectedType]
	if !ok {
		return fmt.Errorf("file type %s is not supported", expectedType)
	}

	signature := make([]byte, len(expectedSignature))
	if _, err := file.Read(signature); err != nil {
		return fmt.Errorf("failed to read file signature: %v", err)
	}

	if !bytes.Equal(signature, expectedSignature) {
		return fmt.Errorf("file signature does not match expected for %s", expectedType)
	}

	// Возвращение указателя в начало файла
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek file: %v", err)
	}

	return nil
}

func validateEmails(emails []string) []string {
	validEmails := make([]string, 0, len(emails))
	for _, email := range emails {
		if EmailRX.MatchString(email) {
			validEmails = append(validEmails, email)
		}
	}
	return validEmails
}
