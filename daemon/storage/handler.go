package storage

//storage/handler.go

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
