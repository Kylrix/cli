package cmd

import (
	"fmt"

	"github.com/nathfavour/kylrix/cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	decrypt bool
)

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage Kylrix Vault secrets",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var vaultListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Banner("Kylrix Vault - Secrets")
		secrets := []string{"github-token", "aws-access-key", "prod-db-password"}
		for _, s := range secrets {
			fmt.Printf("- %s\n", s)
		}
	},
}

var vaultGetCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Retrieve a secret",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		utils.Banner("Kylrix Vault - Get")
		
		if decrypt {
			utils.Info(fmt.Sprintf("Decrypting secret: %s", name))
			fmt.Printf("Value: ******** (Decrypted)\n")
		} else {
			fmt.Printf("Value: [ENCRYPTED]\n")
			utils.Info("Use --decrypt to see the value")
		}
	},
}

func init() {
	vaultGetCmd.Flags().BoolVarP(&decrypt, "decrypt", "d", false, "Decrypt the secret value")
	
	vaultCmd.AddCommand(vaultListCmd)
	vaultCmd.AddCommand(vaultGetCmd)
	rootCmd.AddCommand(vaultCmd)
}
