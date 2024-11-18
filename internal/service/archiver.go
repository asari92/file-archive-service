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

	"github.com/mholt/archiver/v3"
)

type ArchiveService struct {
	Archiver Archiver
}

func NewArchiverService(archiver Archiver) *ArchiveService {
	return &ArchiveService{
		Archiver: archiver,
	}
}

func (s *ArchiveService) CreateArchive(files []*multipart.FileHeader) (*bytes.Buffer, error) {
	return s.Archiver.CreateArchive(files)
}

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

	// Копируем содержимое файла напрямую во временный файл
	if _, err := io.Copy(tempFile, *mFile); err != nil {
		return nil, errors.New("error writing file data to temp file")
	}

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
