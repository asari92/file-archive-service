package archive

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
	"time"
)

// ZipArchiver is an implementation of Archiver for zip files.
type ZipArchiver struct{}

func NewZipArchiver() *ZipArchiver {
	return &ZipArchiver{}
}

func (za *ZipArchiver) CreateArchive(files []*multipart.FileHeader) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}

		zipFile, err := zipWriter.Create(fileHeader.Filename)
		if err != nil {
			file.Close()
			return nil, err
		}

		if _, err := io.Copy(zipFile, file); err != nil {
			file.Close()
			return nil, err
		}
		file.Close()
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

// TarArchiver is an implementation of Archiver for tar files.
type TarArchiver struct{}

func NewTarArchiver() *TarArchiver {
	return &TarArchiver{}
}

func (ta *TarArchiver) CreateArchive(files []*multipart.FileHeader) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	tarWriter := tar.NewWriter(buf)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		// Создание заголовка для файла с базовыми правами доступа и текущим временем
		header := &tar.Header{
			Name:    fileHeader.Filename,
			Size:    fileHeader.Size,
			Mode:    0o644,      // Базовые права доступа (rw-r--r--)
			ModTime: time.Now(), // Используем текущее время
		}

		// Запись заголовка в архив
		if err := tarWriter.WriteHeader(header); err != nil {
			return nil, err
		}

		// Копирование содержимого файла в архив
		if _, err := io.Copy(tarWriter, file); err != nil {
			return nil, err
		}
	}

	// Закрытие tar.Writer для завершения записи архива
	if err := tarWriter.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}
