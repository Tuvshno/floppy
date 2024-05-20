package upload

//upload/handler.go

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Handler struct{}

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
}
