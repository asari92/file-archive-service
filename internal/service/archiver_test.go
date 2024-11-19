package service

import (
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"testing"
)

func TestService_GenerateArchiveInfo(t *testing.T) {
	service := Service{}
	// Создание временного файла, который будет содержать данные архива
	tempFile, err := os.CreateTemp("", "test-*.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name()) // Удалить после теста

	// Создание ZIP архива
	zipWriter := zip.NewWriter(tempFile)
	fileInZip, err := zipWriter.Create("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, err = fileInZip.Write([]byte("content of the zip file"))
	if err != nil {
		t.Fatal(err)
	}
	zipWriter.Close()
	tempFile.Close()

	// Открываем файл архива для чтения
	zipFile, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer zipFile.Close()

	// Создание multipart.File и FileHeader для тестирования
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "testfile.zip")
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, zipFile)
	if err != nil {
		t.Fatal(err)
	}
	writer.Close()

	// Чтение созданной формы
	r := multipart.NewReader(body, writer.Boundary())
	form, err := r.ReadForm(1024)
	if err != nil {
		t.Fatal(err)
	}
	file := form.File["file"][0]
	mFile, err := file.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer mFile.Close()

	// Симуляция загрузки файла и вызов функции
	tests := []struct {
		name    string
		mFile   multipart.File
		mHeader *multipart.FileHeader
		wantErr bool
	}{
		{
			name:    "Valid Input",
			mFile:   mFile,
			mHeader: file,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GenerateArchiveInfo(&tt.mFile, tt.mHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GenerateArchiveInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Errorf("Expected non-nil response, got nil")
			}
		})
	}
}
