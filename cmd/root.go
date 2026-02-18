package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kylrix",
	Short: "Kylrix Ecosystem CLI",
	Long:  `A robust CLI tool for managing the Kylrix ecosystem (Note, Vault, Connect, Keep, Flow, Accounts).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to the Kylrix Ecosystem CLI")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
