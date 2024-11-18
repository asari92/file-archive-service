package service

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"

	"file-archive-service/pkg/config"
	"file-archive-service/pkg/validator"
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
		Archiver: NewArchiveUsecases(archiver),
		Mailer:   NewMailUsecases(mailer),
		Config:   conf,
	}
}

func (s *Service) ProcessAndSendFile(fileHeader *multipart.FileHeader, emails []string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	if fileHeader.Size > s.Config.MaxSendFileSize {
		return fmt.Errorf("file size exceeds the maximum limit")
	}

	if err := validator.ValidateFileSignature(file, fileHeader.Header.Get("Content-Type"), mailerAllowedSignatures); err != nil {
		return err
	}

	recipients, err := validator.ValidateEmails(emails)
	if err != nil {
		return err
	}

	return s.Mailer.SendEmailWithAttachment(s.Config.MailFrom, recipients, "Document", fileHeader.Filename, "", file)
}
