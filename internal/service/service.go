package service

import (
	"mime/multipart"

	"file-archive-service/pkg/config"
)

type Mail interface {
	SendEmailWithAttachment(file multipart.File, filename string, recipients []string) error
}

type Service struct {
	Mail
}

func NewService(conf *config.Config) *Service {
	return &Service{
		Mail: NewMailUsecases(conf),
	}
}
