package cmd

import (
	"fmt"

	"github.com/nathfavour/kylrix/cli/pkg/config"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Kylrix Ecosystem",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Skeleton for login logic
		fmt.Println("Attempting to login...")
		
		// In a real scenario, we'd prompt for credentials or use WebAuthn
		// For now, let's just simulate saving a config
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		
		fmt.Print("Enter API Key: ")
		var apiKey string
		fmt.Scanln(&apiKey)
		
		cfg.APIKey = apiKey
		err = config.SaveConfig(cfg)
		if err != nil {
			return err
		}
		
		fmt.Println("Successfully logged in and saved configuration.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
