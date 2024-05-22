package storage

//storage/handler.go

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tuvshno/floppy/daemon/types"
)

type Handler struct{}

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
	return nil
}
