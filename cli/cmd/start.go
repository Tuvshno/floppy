package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the floppy server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start")
	},
}
