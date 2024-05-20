package main

//router.go

import (
	"fmt"
	"net/http"

	"github.com/tuvshno/floppy/daemon/upload"
)

func loadRoutes(router *http.ServeMux) {
	fmt.Println("Inside loadroutes")
	handler := upload.Handler{}

	router.HandleFunc("POST /upload", handler.Upload)
}

func hello() {
	fmt.Println("hello")
}
