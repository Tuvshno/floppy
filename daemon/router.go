package main

//router.go

import (
	"net/http"

	"github.com/tuvshno/floppy/daemon/upload"
)

// loadRoutes loads the routes from specific handlers to the main servemux multiplexer
func loadRoutes(router *http.ServeMux) {
	handler := upload.Handler{}

	router.HandleFunc("POST /upload", handler.Upload)
}
