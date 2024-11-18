package service

import (
	"bytes"
	"io"
	"mime/multipart"

	"file-archive-service/pkg/config"
)

// Mailer определяет интерфейс для отправки email
type Mailer interface {
	SendEmailWithAttachment(from string, to []string, subject, filename, text string, data io.Reader) error
}

type Archiver interface {
	CreateArchive(files []*multipart.FileHeader) (*bytes.Buffer, error)
	// GenerateArchiveInfo(mFile *multipart.File, mFileHeader *multipart.FileHeader) (*models.Response, error)
}

type Service struct {
	Archiver
	Mailer
	Config *config.Config
}

func NewService(archiver Archiver, mailer Mailer, conf *config.Config) *Service {
	return &Service{
		Archiver: NewArchiverService(archiver),
		Mailer:   NewMailerService(mailer),
		Config:   conf,
	}
}
