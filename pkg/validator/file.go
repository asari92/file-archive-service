package validator

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
)

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
