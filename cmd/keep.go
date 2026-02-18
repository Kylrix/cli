package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var keepCmd = &cobra.Command{
	Use:   "keep",
	Short: "Manage Kylrix Keep",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Kylrix Keep management")
	},
}

func init() {
	rootCmd.AddCommand(keepCmd)
}
