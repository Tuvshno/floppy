package main

//router.go

import (
	"net/http"

	"github.com/tuvshno/floppy/daemon/storage"
	"github.com/tuvshno/floppy/daemon/upload"
)

// loadRoutes loads the routes from specific handlers to the main servemux multiplexer
func loadRoutes(router *http.ServeMux) {
	uploadHandler := upload.Handler{}
	storageHandler := storage.Handler{}

	router.HandleFunc("POST /upload", uploadHandler.Upload)
	router.HandleFunc("POST /storage", storageHandler.Store)
}
