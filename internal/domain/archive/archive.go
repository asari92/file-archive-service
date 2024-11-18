package archive

import (
	// "archive/tar"
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
)

// ZipArchiver is an implementation of Archiver for zip files.
type ZipArchiver struct{}

func NewZipArchiver() *ZipArchiver {
	return &ZipArchiver{}
}

func (s *ZipArchiver) CreateArchive(files []*multipart.FileHeader) (*bytes.Buffer, error) {
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

// func (s *ZipArchiver) GenerateArchiveInfo(mFile *multipart.File, mFileHeader *multipart.FileHeader) (*models.Response, error) {
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

// TarArchiver is an implementation of Archiver for tar files.
type TarArchiver struct{}

func NewTarArchiver() *TarArchiver {
	return &TarArchiver{}
}

// func (t *TarArchiver) CreateArchive(files []*multipart.FileHeader) (*bytes.Buffer, error) {
// 	buf := new(bytes.Buffer)
// 	tarWriter := tar.NewWriter(buf)

// 	for _, fileHeader := range files {
// 		file, err := fileHeader.Open()
// 		if err != nil {
// 			return nil, err
// 		}

// 		tarWriter.Write()
// 		zipFile, err := tarWriter.Create(fileHeader.Filename)
// 		if err != nil {
// 			file.Close()
// 			return nil, err
// 		}

// 		if _, err := io.Copy(zipFile, file); err != nil {
// 			file.Close()
// 			return nil, err
// 		}
// 		file.Close()
// 	}

// 	if err := tarWriter.Close(); err != nil {
// 		return nil, err
// 	}

// 	return buf, nil
// }
