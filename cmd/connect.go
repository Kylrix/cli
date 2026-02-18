package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Manage Kylrix Connect",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Kylrix Connect management")
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
