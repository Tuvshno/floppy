package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run server.go <filename>")
		return
	}

	filename := os.Args[1]

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()
	log.Println("Server Listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
		}
		go handleConnection(conn, filename)
	}
}

func handleConnection(conn net.Conn, filename string) {
	defer conn.Close()

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Failed to open %s", "largefile.txt")
		return
	}
	defer file.Close()

	buffer := make([]byte, 4096)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err.Error() != "EOF" {
				log.Println("Error reading file:", err)
			}
			break
		}

		_, err = conn.Write(buffer[:n])
		if err != nil {
			log.Println("Error sending data:", err)
			break
		}
	}

	fmt.Println("File transfer complete")
}
