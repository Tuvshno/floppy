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
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload [file path]",
	Short: "Uploads a file to the floppy daemon",
	Long:  `Calls to the floppy daemon to upload a file that will be availible to all clients`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userFilePath := args[0]
		fmt.Println("upload called with file", args[0])

		isWSL, err := cmd.Flags().GetBool("wsl")
		if err != nil {
			fmt.Printf("Failed to get flag value: %v\n", err)
			return
		}
		if isWSL {
			wslPath, err := convertToWSLPath(userFilePath)
			if err != nil {
				fmt.Printf("Failed to convert path to WSL path: %v\n", err)
				return
			}
			userFilePath = wslPath
		}

		absPath, err := filepath.Abs(userFilePath)
		if err != nil {
			fmt.Printf("Failed to get absolute path %v\n", err)
			return
		}
		fmt.Println("Absolute file path:", absPath)

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			fmt.Println("File does not exist")
			return
		}

		file, err := os.Open(absPath)
		if err != nil {
			fmt.Println("File does not exist")
			return
		}
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", filepath.Base(absPath))
		if err != nil {
			fmt.Printf("Failed to create form file: %v\n", err)
			return
		}

		_, err = io.Copy(part, file)
		if err != nil {
			fmt.Println("Failed to copy file content %v\n", err)
			return
		}

		err = writer.Close()
		if err != nil {
			fmt.Println("Failed to close writer %v\n", err)
			return
		}

		request, err := http.NewRequest("POST", "http://localhost:8080/upload", body)
		if err != nil {
			fmt.Println("Failed to create request %v\n", err)
		}
		request.Header.Add("Content-Type", writer.FormDataContentType())

		client := &http.Client{
			// Timeout: 30 * time.Second,
		}

		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Failed to execute request %v\n", err)
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			fmt.Println("Failed to upload file, status code %d\n", response.StatusCode)
			return
		}

		respBody, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Failed to read response body %v\n", err)
			return
		}

		fmt.Println("Successfully uploaded file: %s\n", respBody)

	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().Bool("wsl", false, "Specify if running on WSL to convert filepaths")
}

// convertToWSLPath converts a Windows path to a WSL path
func convertToWSLPath(winPath string) (string, error) {
	cmd := exec.Command("wslpath", "-a", winPath)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
