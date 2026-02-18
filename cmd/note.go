package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage Kylrix Note",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Kylrix Note management")
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)
}
