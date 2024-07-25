package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Progress Reader is a struct that implements the io.Reader to show progress of an upload
type ProgressReader struct {
	io.Reader
	Total    int64
	Current  int64
	FileSize int64
}

type Metadata struct {
	FilePath string `json:"file_path"`
}

// Read overrides the io.Read interface to update the current bytes read
func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.Current += int64(n)
	pr.updateProgress()
	return n, err
}

// updateProgress prints the current progress of the upload to the terminal
func (pr *ProgressReader) updateProgress() {
	percentage := float64(pr.Current) / float64(pr.FileSize) * 100
	fmt.Printf("\rUploading... %d/%d bytes (%.2f%%)", pr.Current, pr.FileSize, percentage)
}

// GetLocalIP retrieves the non-loopback, non-link-local IP address of the host machine
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil && !ipNet.IP.IsLinkLocalUnicast() {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", errors.New("unable to determine non-link-local IP address")
}

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload [file path]",
	Short: "Uploads a file to the floppy daemon",
	Long:  `Calls to the floppy daemon to upload a file that will be available to all clients`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userFilePath := args[0]

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

		fileInfo, err := file.Stat()
		if err != nil {
			fmt.Printf("Failed getting file stats %v\n", err)
			return
		}
		fileSize := fileInfo.Size()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", filepath.Base(absPath))
		if err != nil {
			fmt.Printf("Failed to create form file: %v\n", err)
			return
		}

		progressReader := &ProgressReader{
			Reader:   file,
			FileSize: fileSize,
		}

		_, err = io.Copy(part, progressReader)
		if err != nil {
			fmt.Printf("Failed to copy file content %v\n", err)
			return
		}

		metadata := Metadata{FilePath: absPath}
		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			fmt.Printf("Failed to marshal metadata %v\n", err)
			return
		}

		metadataPart, err := writer.CreateFormField("metadata")
		if err != nil {
			fmt.Printf("Failed to create form field for metadata %v\n", err)
			return
		}

		_, err = metadataPart.Write(metadataJSON)
		if err != nil {
			fmt.Printf("Failed to write metadata %v\n", err)
			return
		}

		err = writer.Close()
		if err != nil {
			fmt.Printf("Failed to close writer %v\n", err)
			return
		}

		localIP, err := GetLocalIP()
		if err != nil {
			fmt.Printf("Failed to get local IP: %v\n", err)
			return
		}

		request, err := http.NewRequest("POST", fmt.Sprintf("http://%s:8080/upload", localIP), body)
		if err != nil {
			fmt.Printf("Failed to create request %v\n", err)
			return
		}
		request.Header.Add("Content-Type", writer.FormDataContentType())

		client := &http.Client{}

		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("Failed to execute request %v\n", err)
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			fmt.Printf("Failed to upload file, status code %d\n", response.StatusCode)
			return
		}

		respBody, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Failed to read response body %v\n", err)
			return
		}

		fmt.Printf("\n%s\n", string(respBody))

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
