package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage Kylrix Vault",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Kylrix Vault management")
	},
}

func init() {
	rootCmd.AddCommand(vaultCmd)
}
