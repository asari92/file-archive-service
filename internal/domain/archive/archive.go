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
