package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tuvshno/floppy/cli/network"
)

// listCmd represents the list command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Finds the IP of the Daemon",
	Long:  `Calls to the Service Discovery Module to find the first Remote IP Daemon`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr, err := network.DiscoverService()
		if err != nil {
			fmt.Printf("Failed to discover service: %v\n", err)
			return
		}

		fmt.Println(serverAddr)
	},
}

func init() {
	rootCmd.AddCommand(discoverCmd)
}
