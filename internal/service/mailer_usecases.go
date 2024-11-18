package service

import (
	"io"
)

var allowedSignatures = map[string][]byte{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {0x50, 0x4B, 0x03, 0x04}, // DOCX
	"application/pdf": {0x25, 0x50, 0x44, 0x46}, // PDF
}

type MailerUsecases struct {
	Mailer Mailer
}

func NewMailUsecases(mailer Mailer) *MailerUsecases {
	return &MailerUsecases{Mailer: mailer}
}

func (ms *MailerUsecases) SendEmailWithAttachment(from string, to []string, subject, filename, text string, data io.Reader) error {
	return ms.Mailer.SendEmailWithAttachment(from, to, subject, filename, text, data)
}
