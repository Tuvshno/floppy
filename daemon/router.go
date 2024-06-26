package main

//router.go

import (
	"database/sql"
	"net/http"

	"github.com/tuvshno/floppy/daemon/storage"
	"github.com/tuvshno/floppy/daemon/upload"
)

// loadRoutes loads the routes from specific handlers to the main servemux multiplexer
func loadRoutes(router *http.ServeMux, db *sql.DB) {
	uploadHandler := upload.Handler{}
	storageHandler := storage.Handler{DB: db}

	router.HandleFunc("POST /upload", uploadHandler.Upload)

	router.HandleFunc("GET /storage", storageHandler.List)
	router.HandleFunc("POST /storage", storageHandler.Store)
	router.HandleFunc("DELETE /storage", storageHandler.Delete)
}
