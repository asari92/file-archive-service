package mail

import (
	"io"
	"time"

	"file-archive-service/pkg/config"

	gomail "gopkg.in/mail.v2"
)

// Dialer интерфейс для мокирования отправки сообщений
type Dialer interface {
	DialAndSend(m ...*gomail.Message) error
}

// GoMailAdapter использует интерфейс Dialer для отправки почты
type GoMailAdapter struct {
	Dialer Dialer
}

func NewGoMailAdapter(conf *config.Config) *GoMailAdapter {
	// Инициализация адаптера
	dialer := gomail.NewDialer(conf.SMTPHost, conf.SMTPPort, conf.SMTPUser, conf.SMTPPassword)
	dialer.Timeout = conf.DialerTimeout * time.Second
	return &GoMailAdapter{Dialer: dialer}
}

func (adapter *GoMailAdapter) SendEmailWithAttachment(from string, to []string, subject, filename, text string, data io.Reader) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.AttachReader(filename, data)
	return adapter.Dialer.DialAndSend(m)
}
