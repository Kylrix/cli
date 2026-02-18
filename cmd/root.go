package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "kylrix",
	Short: "Kylrix Ecosystem CLI",
	Long:  `A robust CLI tool for managing the Kylrix ecosystem (Note, Vault, Connect, Keep, Flow, Accounts).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to the Kylrix Ecosystem CLI")
		fmt.Println("Use 'kylrix --help' for available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kylrix/config.json)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
