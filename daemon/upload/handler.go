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

// Metadata holds the struct recieved from the client
type Metadata struct {
	FilePath string `json:"file_path"`
}

// Upload recieves a POST request that uploads the given file inside the file form to the daemon
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

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to copy file to save", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("File Uploaded Successfully"))
	log.Println("Uploaded File Succesfully")

	metadataStr := r.FormValue("metadata")
	if metadataStr == "" {
		http.Error(w, "Metadata not provided", http.StatusBadRequest)
		return
	}

	var clientMetadata Metadata
	err = json.Unmarshal([]byte(metadataStr), &clientMetadata)
	if err != nil {
		http.Error(w, "Failed to unmarshal metadata", http.StatusInternalServerError)
		return
	}

	metadata := types.FileMetadata{
		Filename: filename,
		Size:     header.Size,
		UploadAt: time.Now(),
		FilePath: clientMetadata.FilePath,
	}

	jsonMetadata, err := json.Marshal(metadata)
	if err != nil {
		fmt.Printf("Failed to marshal metadata %v", err)
		return
	}

	request, err := http.NewRequest("POST", "http://localhost:8080/storage", bytes.NewBuffer(jsonMetadata))
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

	log.Printf("%s\n", string(respBody))

}
