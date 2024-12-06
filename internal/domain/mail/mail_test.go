package mail

import (
	"fmt"
	"io"
	"log"
	"strings"
	"testing"

	"file-archive-service/pkg/config"
	"file-archive-service/pkg/utils"

	gomail "gopkg.in/mail.v2"
)

type MockDialer struct {
	SendFunc func(m ...*gomail.Message) error
}

func (md *MockDialer) DialAndSend(m ...*gomail.Message) error {
	if md.SendFunc != nil {
		return md.SendFunc(m...)
	}
	return nil
}

func TestGoMailAdapter_SendEmailWithAttachment(t *testing.T) {
	utils.InitAbsolutePath()
	// Загрузите переменные окружения из файла .env
	if err := utils.LoadEnv(utils.GetAbsPath() + "/.env"); err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}

	conf := config.New()

	fmt.Println(conf)

	tests := []struct {
		name     string
		dialer   *MockDialer
		from     string
		to       []string
		subject  string
		filename string
		text     string
		data     io.Reader
		wantErr  bool
	}{
		{
			name:     "Success - Email sent",
			from:     "sender@example.com",
			to:       []string{"test@gmail.com"},
			subject:  "Test Subject",
			filename: "attachment.txt",
			text:     "Hello, this is a test",
			data:     strings.NewReader("This is the content of the attachment"),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewGoMailAdapter(conf)
			err := adapter.SendEmailWithAttachment(tt.from, tt.to, tt.subject, tt.filename, tt.text, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoMailAdapter.SendEmailWithAttachment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
