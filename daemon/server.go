package main

import (
	"github.com/grandcat/zeroconf"
	"github.com/tuvshno/floppy/daemon/storage"
	"log"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	log.Println("Received Request")
	w.Write([]byte("Hello from daemon\n"))
}

func main() {
	// Start the mDNS service
	server, err := zeroconf.Register("FloppyDaemon", "_http._tcp", "local.", 8080, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		log.Fatalf("Failed to register mDNS service: %v", err)
	}
	defer server.Shutdown()

	router := http.NewServeMux()
	router.HandleFunc("/", handle)

	db, err := storage.InitDB("files.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Print("Initiated DB")
	loadRoutes(router, db)

	serverHTTP := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Server Listening on Port :8080")
	log.Fatal(serverHTTP.ListenAndServe())
}
