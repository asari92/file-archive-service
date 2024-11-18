package mail

import (
	"io"
	"time"

	gomail "gopkg.in/mail.v2"
)

type GoMailAdapter struct {
	Dialer *gomail.Dialer
}

func NewGoMailAdapter(host string, port int, username, password string, timeout time.Duration) *GoMailAdapter {
	dialer := gomail.NewDialer(host, port, username, password)
	dialer.Timeout = timeout
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
