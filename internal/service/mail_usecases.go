package service

import (
	"mime/multipart"

	"file-archive-service/pkg/config"

	gomail "gopkg.in/mail.v2"
)

type MailUsecases struct {
	dialer *gomail.Dialer
}

func NewMailUsecases(conf *config.Config) *MailUsecases {
	return &MailUsecases{
		dialer: &gomail.Dialer{
			Host: conf.SMTPHost, Port: conf.SMTPPort,
			Username: conf.SMTPUser, Password: conf.SMTPPassword, Timeout: conf.DialerTimeout,
		},
	}
}

func (ms *MailUsecases) SendEmailWithAttachment(file multipart.File, filename string, recipients []string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", ms.dialer.Username)
	message.SetHeader("To", recipients...)
	message.SetHeader("Subject", "Document")
	message.AttachReader(filename, file)
	return ms.dialer.DialAndSend(message)
}
