package cmd

import (
	"github.com/nathfavour/kylrix/cli/pkg/config"
	"github.com/nathfavour/kylrix/cli/pkg/utils"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Kylrix Ecosystem",
	RunE: func(cmd *cobra.Command, args []string) error {
		utils.Banner("Kylrix Authentication")
		
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		baseURI, err := utils.Prompt("API Base URI (default: https://api.kylrix.com)")
		if err != nil {
			return err
		}
		if baseURI != "" {
			cfg.BaseURI = baseURI
		}

		apiKey, err := utils.PasswordPrompt("Enter API Key")
		if err != nil {
			return err
		}
		
		cfg.APIKey = apiKey
		err = config.SaveConfig(cfg)
		if err != nil {
			return err
		}
		
		utils.Success("Successfully authenticated and configuration saved.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
