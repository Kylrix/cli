package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage Kylrix Accounts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Kylrix Accounts management")
	},
}

func init() {
	rootCmd.AddCommand(accountsCmd)
}
