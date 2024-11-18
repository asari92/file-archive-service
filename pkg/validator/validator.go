package validator

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"regexp"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func ValidateEmails(emails []string) ([]string, error) {
	validEmails := make([]string, 0, len(emails))
	for _, email := range emails {
		if EmailRX.MatchString(email) {
			validEmails = append(validEmails, email)
		}
	}
	if len(validEmails) == 0 {
		return nil, fmt.Errorf("no valid emails provided")
	}

	return validEmails, nil
}

func ValidateFileSignature(file multipart.File, expectedType string, allowedSignatures map[string][]byte) error {
	expectedSignature, ok := allowedSignatures[expectedType]
	if !ok {
		return fmt.Errorf("file type %s is not supported", expectedType)
	}

	signature := make([]byte, len(expectedSignature))
	if _, err := file.Read(signature); err != nil {
		return fmt.Errorf("failed to read file signature: %v", err)
	}

	if !bytes.Equal(signature, expectedSignature) {
		return fmt.Errorf("file signature does not match expected for %s", expectedType)
	}

	// Возвращение указателя в начало файла
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek file: %v", err)
	}

	return nil
}

func ValidateMimeTypes(files []*multipart.FileHeader, archiverAllowedMimeTypes map[string]bool) error {
	for _, fileHeader := range files {
		// Проверяем MIME-тип
		if !archiverAllowedMimeTypes[fileHeader.Header.Get("Content-Type")] {
			return fmt.Errorf("invalid file type: %s", fileHeader.Filename)
		}
	}

	return nil
}
