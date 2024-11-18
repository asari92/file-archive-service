package service

import (
	"io"
)

type MailerUsecases struct {
	Mailer Mailer
}

func NewMailUsecases(mailer Mailer) *MailerUsecases {
	return &MailerUsecases{Mailer: mailer}
}

func (ms *MailerUsecases) SendEmailWithAttachment(from string, to []string, subject, filename, text string, data io.Reader) error {
	return ms.Mailer.SendEmailWithAttachment(from, to, subject, filename, text, data)
}
