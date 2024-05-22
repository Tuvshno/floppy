package types

import "time"

//types/structs.go

type FileMetadata struct {
	Filename string    `json:"filename"`
	Size     int64     `json:"size"`
	UploadAt time.Time `json:"upload_at"`
	FilePath string    `json:"file_path"`
}
