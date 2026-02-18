package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var flowCmd = &cobra.Command{
	Use:   "flow",
	Short: "Manage Kylrix Flow",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Kylrix Flow management")
	},
}

func init() {
	rootCmd.AddCommand(flowCmd)
}
