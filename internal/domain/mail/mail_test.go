package mail

import (
	"errors"
	"io"
	"strings"
	"testing"

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
			name: "Success - Email sent",
			dialer: &MockDialer{
				SendFunc: func(m ...*gomail.Message) error {
					return nil // No error, simulating successful send
				},
			},
			from:     "sender@example.com",
			to:       []string{"receiver@example.com"},
			subject:  "Test Subject",
			filename: "attachment.txt",
			text:     "Hello, this is a test",
			data:     strings.NewReader("This is the content of the attachment"),
			wantErr:  false,
		},
		{
			name: "Failure - SMTP connection error",
			dialer: &MockDialer{
				SendFunc: func(m ...*gomail.Message) error {
					return errors.New("SMTP connection error")
				},
			},
			from:     "sender@example.com",
			to:       []string{"receiver@example.com"},
			subject:  "Test Subject",
			filename: "attachment.txt",
			text:     "Hello, this is a test",
			data:     strings.NewReader("This is the content of the attachment"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewGoMailAdapter(tt.dialer)
			err := adapter.SendEmailWithAttachment(tt.from, tt.to, tt.subject, tt.filename, tt.text, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoMailAdapter.SendEmailWithAttachment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
