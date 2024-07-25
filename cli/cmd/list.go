package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/tuvshno/floppy/cli/network"
	"io"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/tuvshno/floppy/cli/types"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the current files stored in the daemon",
	Long:  `Calls to the floppy daemon to list the current files stored`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr, err := network.DiscoverService()
		if err != nil {
			fmt.Printf("Failed to discover service: %v\n", err)
			return
		}

		request, err := http.NewRequest("GET", fmt.Sprintf("http://%s/storage", serverAddr), nil)
		if err != nil {
			fmt.Printf("Failed to create request %v\n", err)
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("Failed to execute request %v\n", err)
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			fmt.Printf("Failed to get files, status code %d\n", response.StatusCode)
			return
		}

		respBody, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Failed to read response body %v\n", err)
			return
		}

		var files []types.FileMetadata
		err = json.Unmarshal(respBody, &files)
		if err != nil {
			fmt.Printf("Failed to parse response body %v\n", err)
			return
		}

		displayFiles(files)

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func displayFiles(files []types.FileMetadata) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintln(w, "ID\tFilename\tSize\tUpload At\tFile Path")

	for _, file := range files {
		fmt.Fprintf(w, "%d\t%s\t%d\t%s\t%s\n", file.ID, file.Filename, file.Size, file.UploadAt, file.FilePath)
	}

	if len(files) == 0 {
		fmt.Fprint(w, "No files currently in daemon\n")

	}
}
