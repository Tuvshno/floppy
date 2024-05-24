package storage

//storage/handler.go

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("Recieved Delete Request ", http.MethodDelete)

	var metadata types.FileMetadata
	err := json.NewDecoder(r.Body).Decode(&metadata)
	if err != nil {
		http.Error(w, "Failed to decode metadata", http.StatusInternalServerError)
		return
	}

	query := "DELETE FROM files WHERE "
	args := []interface{}{}
	conditions := []string{}

	if metadata.ID != 0 {
		conditions = append(conditions, "ID = ?")
		args = append(args, metadata.ID)
	}
	if metadata.Filename != "" {
		conditions = append(conditions, "Filename = ?")
		args = append(args, metadata.Filename)
	}
	if metadata.Size != 0 {
		conditions = append(conditions, "Size = ?")
		args = append(args, metadata.Size)
	}
	if !metadata.UploadAt.IsZero() {
		if metadata.UploadAt.Hour() == 0 && metadata.UploadAt.Minute() == 0 && metadata.UploadAt.Second() == 0 {
			startOfDay := metadata.UploadAt.Format("2006-01-02 15:04:05 -0700 EDT")
			endOfDay := metadata.UploadAt.Add(24 * time.Hour).Format("2006-01-02 15:04:05 -0700 EDT")
			conditions = append(conditions, "upload_at >= ? AND upload_at < ?")
			args = append(args, startOfDay, endOfDay)
		} else {
			conditions = append(conditions, "upload_at = ?")
			args = append(args, metadata.UploadAt.Format("2006-01-02 15:04:05 -0700 EDT"))
		}
	}
	if metadata.FilePath != "" {
		conditions = append(conditions, "file_path = ?")
		args = append(args, metadata.FilePath)
	}

	if len(conditions) == 0 {
		http.Error(w, "No conditions were given", http.StatusBadRequest)
		return
	}

	query += strings.Join(conditions, " AND ")

	log.Printf("%s %s \n", query, args)

	result, err := h.DB.Exec(query, args...)
	if err != nil {
		http.Error(w, "Failed to execute query", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve rows affected", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No Files matched given flags"))
	} else {
		w.Write([]byte(fmt.Sprintf("Successfully delete %d file(s)", rowsAffected)))
	}

}

func (h *Handler) saveMetadata(metadata types.FileMetadata) error {
	_, err := h.DB.Exec(`
		INSERT INTO files (filename, size, upload_at, file_path)
		VALUES (?,?,?,?)`,
		metadata.Filename, metadata.Size, metadata.UploadAt, metadata.FilePath)
	if err != nil {
		return fmt.Errorf("failed to save metadata: %v", err)
	}
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
