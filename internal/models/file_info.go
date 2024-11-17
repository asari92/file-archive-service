package models

// FileInfo - информация о файле внутри архива
type FileInfo struct {
	FilePath string  `json:"file_path"`
	Size     float64 `json:"size"`
	MimeType string  `json:"mimetype"`
}
