package cmd

import (
	"fmt"

	"github.com/nathfavour/kylrix/cli/pkg/config"
	"github.com/nathfavour/kylrix/cli/pkg/utils"
	"github.com/spf13/cobra"
)

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage Kylrix Accounts and sessions",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var accountsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check current account status",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Banner("Kylrix Accounts - Status")
		cfg, err := config.LoadConfig()
		if err != nil {
			utils.Error("Failed to load configuration")
			return
		}
		
		if cfg.APIKey == "" && cfg.Token == "" {
			utils.Warning("Not logged in. Use 'kylrix login' to authenticate.")
			return
		}
		
		utils.Info(fmt.Sprintf("Logged in to: %s", cfg.BaseURI))
		if cfg.APIKey != "" {
			utils.Info("API Key: [CONFIGURED]")
		}
	},
}

var accountsLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear local authentication data",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Banner("Kylrix Accounts - Logout")
		cfg, err := config.LoadConfig()
		if err != nil {
			utils.Error("Failed to load configuration")
			return
		}
		
		cfg.APIKey = ""
		cfg.Token = ""
		err = config.SaveConfig(cfg)
		if err != nil {
			utils.Error("Failed to clear configuration")
			return
		}
		
		utils.Success("Logged out successfully.")
	},
}

func init() {
	accountsCmd.AddCommand(accountsStatusCmd)
	accountsCmd.AddCommand(accountsLogoutCmd)
	rootCmd.AddCommand(accountsCmd)
}
