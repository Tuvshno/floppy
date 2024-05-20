package upload

//upload/handler.go

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Handler struct{}

type ProgressReader struct {
	io.Reader
	Current  int64
	FileSize int64
}

// Read overides the io.Read interface to update the current bytes read
func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.Current += int64(n)
	pr.updateProgress()
	return n, err
}

// updateProgress pritns the current progress of the upload to the terminal
func (pr *ProgressReader) updateProgress() {
	percentage := float64(pr.Current) / float64(pr.FileSize) * 100
	fmt.Printf("\rUploading... %d/%d bytes (%.2f%%)", pr.Current, pr.FileSize, percentage)
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	log.Println("Recieved Upload Request ", r.Method)
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read form file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filename := header.Filename
	ext := filepath.Ext(filename)

	out, err := os.Create("uploaded_file" + ext)
	if err != nil {
		http.Error(w, "Failed to save new file", http.StatusInternalServerError)
		return
	}

	progressReader := &ProgressReader{
		Reader:   file,
		FileSize: header.Size,
	}

	_, err = io.Copy(out, progressReader)
	if err != nil {
		http.Error(w, "Failed to copy file to save", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("File Uploaded Successfully"))
	fmt.Println("\nFile Uploaded")
}
