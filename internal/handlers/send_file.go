package handlers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	gomail "gopkg.in/mail.v2"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

var allowedMimeTypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/pdf": true,
}

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
	file, err := fileHeader.Open()
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// if err := checkFileType(file, fileHeader.Header.Get("Content-Type")); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	if err := checkFileSignature(file, fileHeader.Header.Get("Content-Type")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получение списка email-адресов из формы
	// emailList := strings.Split(r.FormValue("emails"), ",")
	// if len(emailList) == 0 || emailList[0] == "" {
	// 	http.Error(w, "No valid emails provided", http.StatusBadRequest)
	// 	return
	// }

	// Получение и фильтрация email-адресов
	emailList := strings.Split(r.FormValue("emails"), ",")
	recipients := validateEmails(emailList)
	if len(recipients) == 0 {
		http.Error(w, "No valid emails provided", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	throttle := time.Tick(1 * time.Second) // Ограничение 1 письмо в 200 милисекунд

	dialer := gomail.NewDialer(h.Config.SMTPHost, h.Config.SMTPPort, h.Config.SMTPUser, h.Config.SMTPPassword)

	for _, recipient := range recipients {
		wg.Add(1)
		go sendEmailWithAttachment(file, fileHeader.Filename, recipient, dialer, &wg, throttle)
	}

	wg.Wait() // Ожидаем завершения всех горутин

	// Подготовка и отправка писем
	// if err := h.sendEmailWithAttachment(file, fileHeader.Filename, recipients); err != nil {
	// 	http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
}

// func (h *Handler) sendEmailWithAttachment(file multipart.File, filename string, recipients []string) error {
// 	d := gomail.NewDialer(h.Config.SMTPHost, h.Config.SMTPPort, h.Config.SMTPUser, h.Config.SMTPPassword)

// 	m := gomail.NewMessage()
// 	// m.SetHeader("From", h.Config.SMTPUser)
// 	m.SetHeader("From", "dulmaev.andrei@gmail.com")
// 	m.SetHeader("To", recipients...)
// 	m.SetHeader("Subject", "Document")
// 	m.AttachReader(filename, file)

// 	go func() {
// 		if err := d.DialAndSend(m); err != nil {
// 			fmt.Println("Failed to send email:", err)
// 		}
// 	}()

// 	return nil
// }

func sendEmailWithAttachment(file multipart.File, filename string, recipient string, dialer *gomail.Dialer, wg *sync.WaitGroup, throttle <-chan time.Time) {
	defer wg.Done() // Указываем, что горутина завершена

	<-throttle // Ожидаем разрешение от канала throttle

	message := gomail.NewMessage()
	message.SetHeader("From", "dulmaev.andrei@gmail.com")
	message.SetHeader("To", recipient)
	message.SetHeader("Subject", "Document")
	message.AttachReader(filename, file) // Вложение файла

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Printf("Failed to send email to %s: %v\n", recipient, err)
	}
}

// func (h *Handler) sendEmailWithAttachment(file multipart.File, filename string, recipients []string, contentType string) error {
// 	// Настройка информации для аутентификации
// 	auth := smtp.PlainAuth("", h.Config.SMTPUser, h.Config.SMTPPassword, h.Config.SMTPHost)

// 	// Создание буфера для письма
// 	var email bytes.Buffer
// 	writer := multipart.NewWriter(&email)
// 	defer writer.Close()

// 	// Добавление заголовка
// 	headers := make(map[string]string)
// 	headers["From"] = "github.com/asari92"
// 	headers["To"] = strings.Join(recipients, ",")
// 	headers["Subject"] = "Document"
// 	headers["MIME-Version"] = "1.0"
// 	headers["Content-Type"] = "multipart/mixed; boundary=" + writer.Boundary()
// 	writeHeaders(&email, headers)

// 	// Добавление файла как вложения
// 	part, err := writer.CreatePart(textproto.MIMEHeader{
// 		"Content-Disposition":       {"attachment; filename=\"" + filename + "\""},
// 		"Content-Type":              {contentType},
// 		"Content-Transfer-Encoding": {"base64"},
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	encoder := base64.NewEncoder(base64.StdEncoding, part)
// 	if _, err = io.Copy(encoder, file); err != nil {
// 		return err
// 	}
// 	encoder.Close()

// 	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!TODO  переделать с использование горутин и waitgroup чтобы каждому адресату отдельно обрабатывалось отправка письма
// 	// Отправка письма
// 	return smtp.SendMail(h.Config.SMTPHost+":"+h.Config.SMTPPort, auth, "github.com/asari92", recipients, email.Bytes())
// }

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
