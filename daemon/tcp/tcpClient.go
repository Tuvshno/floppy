package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func startTCPClient() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run client.go <outputfilename>")
		return
	}

	outputFilename := os.Args[1]

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Println("Error connecting to server:", err)
		return
	}

	defer conn.Close()

	file, err := os.Create(outputFilename)
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error receiving data:", err)
			return
		}
		fmt.Println("Recieved: ", n)
		_, err = file.Write(buffer[:n])
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	fmt.Println("File received successfully")

}
