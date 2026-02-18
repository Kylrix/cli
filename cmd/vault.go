package cmd

import (
	"fmt"

	"github.com/nathfavour/kylrix/cli/pkg/crypto"
	"github.com/nathfavour/kylrix/cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	decryptSecret bool
	vaultPassword string
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
		header := []string{"NAME", "CREATED", "TYPE"}
		data := [][]string{
			{"github-token", "2024-02-10", "Static"},
			{"aws-access-key", "2024-02-12", "Static"},
			{"prod-db-password", "2024-02-15", "Encrypted"},
		}
		utils.Table(header, data)
	},
}

var vaultCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new secret",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		
		value, err := utils.PasswordPrompt("Secret Value")
		if err != nil {
			return err
		}

		password, err := utils.PasswordPrompt("Vault Master Password (for key derivation)")
		if err != nil {
			return err
		}

		// Use a dummy salt for now, in reality this would be retrieved or generated
		salt := make([]byte, crypto.SaltSize)
		copy(salt, []byte("kylrix-ecosystem-default-salt-!!")) 

		utils.Info("Deriving key using PBKDF2 (600,000 iterations)...")
		key := crypto.DeriveKey(password, salt)
		
		encrypted, err := crypto.Encrypt(value, key)
		if err != nil {
			return err
		}

		utils.Success(fmt.Sprintf("Secret '%s' encrypted successfully using WESP.", name))
		fmt.Printf("Encrypted Payload: %s\n", encrypted)
		utils.Info("Note: This matches the 'vault/lib/ecosystem/security.ts' implementation.")
		return nil
	},
}

var vaultGetCmd = &cobra.Command{
	Use:   "get [name] [payload]",
	Short: "Retrieve/Decrypt a secret",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		payload := args[1]
		
		utils.Banner("Kylrix Vault - Get")
		
		if decryptSecret {
			password, err := utils.PasswordPrompt("Vault Master Password")
			if err != nil {
				return err
			}

			salt := make([]byte, crypto.SaltSize)
			copy(salt, []byte("kylrix-ecosystem-default-salt-!!")) 

			key := crypto.DeriveKey(password, salt)
			
			decrypted, err := crypto.Decrypt(payload, key)
			if err != nil {
				utils.Error("Decryption failed. Ensure the password and payload are correct.")
				return err
			}

			utils.Success(fmt.Sprintf("Secret '%s' decrypted:", name))
			fmt.Printf("Value: %v\n", decrypted)
		} else {
			fmt.Printf("Name: %s\nPayload: %s\n", name, payload)
			utils.Info("Use --decrypt to see the value")
		}
		return nil
	},
}

func init() {
	vaultGetCmd.Flags().BoolVarP(&decryptSecret, "decrypt", "d", false, "Decrypt the secret value")
	
	vaultCmd.AddCommand(vaultListCmd)
	vaultCmd.AddCommand(vaultGetCmd)
	vaultCmd.AddCommand(vaultCreateCmd)
	rootCmd.AddCommand(vaultCmd)
}
