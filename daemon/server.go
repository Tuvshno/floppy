package main

//server.go

import (
	"log"
	"net/http"

	"github.com/tuvshno/floppy/daemon/storage"
)

func handle(w http.ResponseWriter, r *http.Request) {
	log.Println("Recieved Request")
	w.Write([]byte("Hello from domain\n"))
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", handle)

	db, err := storage.InitDB("files.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Print("Initiated DB")
	loadRoutes(router, db)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Server Listening on Port :8080")
	log.Fatal(server.ListenAndServe())
}
