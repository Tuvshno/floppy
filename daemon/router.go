package main

import "net/http"

func loadRoutes(router *http.ServeMux) {

	handler := &upload.Handler{}

	router.HandleFunc("POST /upload", handler.Upload)
}
