package service

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"file-archive-service/internal/domain/models"
	"file-archive-service/pkg/utils"
	"file-archive-service/pkg/validator"

	"github.com/mholt/archiver/v3"
)

var archiverAllowedMimeTypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/xml": true,
	"image/jpeg":      true,
	"image/png":       true,
}

type ArchiveUsecases struct {
	Archiver Archiver
}

func NewArchiveUsecases(archiver Archiver) *ArchiveUsecases {
	return &ArchiveUsecases{
		Archiver: archiver,
	}
}

func (s *ArchiveUsecases) CreateArchive(files []*multipart.FileHeader) (*bytes.Buffer, error) {
	err := validator.ValidateMimeTypes(files, archiverAllowedMimeTypes)
	if err != nil {
		return nil, err
	}

	return s.Archiver.CreateArchive(files)
}

// func (s *ArchiveUsecases) GenerateArchiveInfo(mFile *multipart.File, mFileHeader *multipart.FileHeader) (*models.Response, error) {
// 	if mFile == nil || mFileHeader == nil {
// 		return nil, errors.New("invalid file or file header")
// 	}

// 	zipReader, err := zip.NewReader(*mFile, mFileHeader.Size)
// 	if err != nil {
// 		return nil, errors.New("file is not a valid zip archive")
// 	}

// 	var files []models.FileInfo
// 	var totalSize int64
// 	for _, f := range zipReader.File {
// 		info := f.FileInfo()
// 		mimeType := utils.GetMimeType(f.Name)
// 		files = append(files, models.FileInfo{
// 			FilePath: f.Name,
// 			Size:     float64(info.Size()),
// 			MimeType: mimeType,
// 		})
// 		totalSize += info.Size()
// 	}

// 	return &models.Response{
// 		Filename:    mFileHeader.Filename,
// 		ArchiveSize: float64(mFileHeader.Size),
// 		TotalSize:   float64(totalSize),
// 		TotalFiles:  len(files),
// 		Files:       files,
// 	}, nil
// }

func (s *Service) GenerateArchiveInfo(mFile *multipart.File, mFileHeader *multipart.FileHeader) (*models.Response, error) {
	if mFile == nil || mFileHeader == nil {
		return nil, errors.New("invalid file or file header")
	}

	// Сохраняем файл временно
	tempFile, err := os.CreateTemp("", "*"+filepath.Ext(mFileHeader.Filename))
	if err != nil {
		return nil, errors.New("error creating a temp file")
	}
	defer os.Remove(tempFile.Name()) // Удалить файл после завершения обработки

	// Копируем содержимое во временный файл
	fileBytes, err := io.ReadAll(*mFile)
	if err != nil {
		return nil, errors.New("error reading file data")
	}
	tempFile.Write(fileBytes)

	// Определяем тип архива и читаем его
	_, err = archiver.ByExtension(tempFile.Name())
	if err != nil {
		return nil, errors.New("file is not a supported archive type")
	}

	// Получаем информацию из архива
	var files []models.FileInfo
	var totalSize int64
	err = archiver.Walk(tempFile.Name(), func(f archiver.File) error {
		info := f.FileInfo
		mimeType := utils.GetMimeType(f.Name())
		files = append(files, models.FileInfo{
			FilePath: f.Name(),
			Size:     float64(info.Size()),
			MimeType: mimeType,
		})
		totalSize += info.Size()
		return nil
	})
	if err != nil {
		return nil, errors.New("error extracting archive")
	}

	return &models.Response{
		Filename:    mFileHeader.Filename,
		ArchiveSize: float64(mFileHeader.Size),
		TotalSize:   float64(totalSize),
		TotalFiles:  len(files),
		Files:       files,
	}, nil
}
