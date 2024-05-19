/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Uploads a file to the floppy daemon",
	Long:  `Calls to the floppy daemon to upload a file that will be availible to all clients`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("upload called with file", args[0])
		filePath := args[0]
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("File does not exist")
			return
		}
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("filename", filePath)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Something went wrong with writer createformfile")
			return
		}

		_, err = io.Copy(part, file)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Something went wrong with copying the file")
			return
		}

		err = writer.Close()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Something went wrong with closing the writer")
			return
		}

		request, err := http.NewRequest("POST", "localhost:8080/upload", body)
		request.Header.Add("Content-Type",
			writer.FormDataContentType())

		client := &http.Client{}
		response, err := client.Do(request)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Something went wrong in the request")
			return
		}
		defer response.Body.Close()

		fmt.Println("Successfully uploaded")
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
