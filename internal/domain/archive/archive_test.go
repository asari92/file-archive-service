package archive

import (
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
	"testing"
)

// TestCreateArchive проверяет корректность создания zip-архива из списка файлов.
func TestCreateArchive(t *testing.T) {
	// Создаем тестовые файлы в памяти
	fileHeaders := []*multipart.FileHeader{
		createTestFileHeader(t, "testfile1.txt", "Hello, world!"),
		createTestFileHeader(t, "testfile2.txt", "Goodbye, world!"),
	}

	// Инициализация ZipArchiver
	archiver := NewZipArchiver()

	// Создание архива
	buffer, err := archiver.CreateArchive(fileHeaders)
	if err != nil {
		t.Fatalf("Failed to create archive: %v", err)
	}

	// Проверка содержимого архива
	checkArchiveContents(t, buffer, map[string]string{
		"testfile1.txt": "Hello, world!",
		"testfile2.txt": "Goodbye, world!",
	})
}

// createTestFileHeader создает FileHeader для тестирования.
func createTestFileHeader(t *testing.T, filename, content string) *multipart.FileHeader {
	t.Helper()
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte(content))
	w.Close()

	// Считывание формы для создания FileHeader
	r := multipart.NewReader(&b, w.Boundary())
	form, err := r.ReadForm(1024)
	if err != nil {
		t.Fatal(err)
	}
	return form.File["file"][0]
}

// checkArchiveContents проверяет, что содержимое архива соответствует ожидаемому.
func checkArchiveContents(t *testing.T, buf *bytes.Buffer, expectedContents map[string]string) {
	t.Helper()
	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range r.File {
		expected, ok := expectedContents[f.Name]
		if !ok {
			t.Errorf("Unexpected file %s in archive", f.Name)
			continue
		}

		rc, err := f.Open()
		if err != nil {
			t.Error(err)
			continue
		}
		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Error(err)
			continue
		}

		if string(content) != expected {
			t.Errorf("Content mismatch for %s: got %s, want %s", f.Name, content, expected)
		}
	}
}
