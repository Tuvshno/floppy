package upload

//upload/handler.go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/tuvshno/floppy/daemon/types"
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

	metadata := types.FileMetadata{
		Filename: filename,
		Size:     header.Size,
		UploadAt: time.Now(),
		FilePath: "uploaded_file" + ext,
	}

	jsonMetadata, err := json.Marshal(metadata)
	if err != nil {
		fmt.Printf("Failed to marshal metadata %v", err)
		return
	}

	request, err := http.NewRequest("POST", "http://localhost:8080/storage", bytes.NewBuffer(jsonMetadata))
	request.Header.Set("X-Custom-Header", "myvalue")
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Printf("Failed to create storage request %v\n", err)
		return
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Failed to execute request %v\n", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Failed to execute Storage Save %d\n", response.StatusCode)
		return
	}

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed to read response body %v\n", err)
		return
	}

	fmt.Printf("\nSuccessfully uploaded file: %s\n", string(respBody))

}
