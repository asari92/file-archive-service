package models

// Response - структура ответа сервера
type Response struct {
	Filename    string     `json:"filename"`
	ArchiveSize float64    `json:"archive_size"`
	TotalSize   float64    `json:"total_size"`
	TotalFiles  int        `json:"total_files"`
	Files       []FileInfo `json:"files"`
}
