package main

import (
	"fmt"
	"log"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	log.Println("Recieved Request")
	w.Write([]byte("Hello from domain\n"))
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", handle)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server Listening on Port :8080")
	server.ListenAndServe()
}
