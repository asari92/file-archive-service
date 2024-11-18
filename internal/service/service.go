package service

import (
	"io"
)

// Mailer определяет интерфейс для отправки email
type Mailer interface {
	SendEmailWithAttachment(from string, to []string, subject, filename, text string, data io.Reader) error
}

type Service struct {
	Mailer
}

func NewService(mailer Mailer) *Service {
	return &Service{
		Mailer: NewMailUsecases(mailer),
	}
}
