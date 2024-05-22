package storage

//storage/handler.go

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tuvshno/floppy/daemon/types"
)

type Handler struct {
	DB *sql.DB
}

func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
	log.Println("Recieved Store Request ", r.Method)

	var metadata types.FileMetadata
	err := json.NewDecoder(r.Body).Decode(&metadata)
	if err != nil {
		http.Error(w, "Failed to decode metadata", http.StatusBadRequest)
		return
	}

	err = h.saveMetadata(metadata)
	if err != nil {
		http.Error(w, "Failed to save metadata", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Metadata stored successfully"))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	log.Println("Recieved List Request ", r.Method)

	files, err := h.getMetadata()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func (h *Handler) saveMetadata(metadata types.FileMetadata) error {
	fmt.Println(metadata)
	_, err := h.DB.Exec(`
		INSERT INTO files (filename, size, upload_at, file_path)
		VALUES (?,?,?,?)`,
		metadata.Filename, metadata.Size, metadata.UploadAt, metadata.FilePath)
	if err != nil {
		return fmt.Errorf("failed to save metadata: %v", err)
	}
	fmt.Println("Sucessfully Stored Metadata")
	return nil
}

func (h *Handler) getMetadata() ([]types.FileMetadata, error) {
	rows, err := h.DB.Query(`SELECT id, filename, size, upload_at, file_path FROM files`)
	if err != nil {
		return nil, fmt.Errorf("failed to save metadata: %v", err)
	}
	defer rows.Close()

	var files []types.FileMetadata
	for rows.Next() {
		var file types.FileMetadata
		var uploadAt string

		err := rows.Scan(&file.ID, &file.Filename, &file.Size, &uploadAt, &file.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		file.UploadAt, err = time.Parse(time.RFC3339, uploadAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse upload time: %v", err)
		}

		files = append(files, file)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return files, nil
}
