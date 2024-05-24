/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/tuvshno/floppy/cli/types"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes file/files from the daemon",
	Long:  `Delete a file or multiple files from the daemon using the ID, File Name, or File Path`,
	Run: func(cmd *cobra.Command, args []string) {

		id, _ := cmd.Flags().GetInt64("id")
		filename, _ := cmd.Flags().GetString("name")
		size, _ := cmd.Flags().GetInt64("size")
		date, _ := cmd.Flags().GetString("date")
		filepath, _ := cmd.Flags().GetString("path")

		if id == 0 && filename == "" && size == 0 && date == "" && filepath == "" {
			fmt.Println("You must provide at least one flag for --id, --filename, --size, --date, or --filepath")
			cmd.Usage()
			return
		}

		var uploadAt time.Time
		var err error
		if date != "" {
			dateFormats := []string{
				"2006-01-02",
				"2006-01-02 15:04:05 -0700 MST",
			}

			for _, format := range dateFormats {
				uploadAt, err = time.Parse(format, date)
				if err == nil {
					break
				}
			}

			if err != nil {
				fmt.Printf("Failed to parse time %v\n", err)
				return
			}
		}

		metadata := types.FileMetadata{
			ID:       id,
			Filename: filename,
			Size:     size,
			UploadAt: uploadAt,
			FilePath: filepath,
		}

		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			fmt.Printf("Failed to marshal metadata %v\n", err)
			return
		}

		request, err := http.NewRequest("DELETE", "http://localhost:8080/storage", bytes.NewBuffer(metadataJSON))
		if err != nil {
			fmt.Printf("Failed to create request %v\n", err)
			return
		}

		client := http.Client{}
		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("Failed to execute request: %v\n", err)
			return
		}
		defer response.Body.Close()

		respBody, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Failed to read body %v\n", err)
			return
		}

		fmt.Println(string(respBody))

		if response.StatusCode != http.StatusOK {
			fmt.Printf("Failed to delete file : %v\n", response.StatusCode)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().Int64("id", 0, "id of the file you want to delete")
	deleteCmd.Flags().String("name", "", "filename of the file you want to delete")
	deleteCmd.Flags().Int64("size", 0, "size of the file you want to delete")
	deleteCmd.Flags().String("date", "", "date of the file you want to delete (YYYY-MM-DD)")
	deleteCmd.Flags().String("path", "", "filepath of the file you want to delete")
}
