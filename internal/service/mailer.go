package service

import (
	"io"
)

type MailerService struct {
	Mailer            Mailer
	AllowedSignatures map[string][]byte
}

func NewMailerService(mailer Mailer) *MailerService {
	return &MailerService{Mailer: mailer}
}

func (ms *MailerService) SendEmailWithAttachment(from string, to []string, subject, filename, text string, data io.Reader) error {
	return ms.Mailer.SendEmailWithAttachment(from, to, subject, filename, text, data)
}
