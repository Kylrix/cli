package cmd

import (
	"fmt"

	"github.com/nathfavour/kylrix/cli/pkg/crypto"
	"github.com/nathfavour/kylrix/cli/pkg/db"
	"github.com/nathfavour/kylrix/cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	decryptSecret bool
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
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.InitDB()
		if err != nil {
			return err
		}
		defer database.Close()

		rows, err := database.Query("SELECT name, created_at FROM vault_secrets")
		if err != nil {
			return err
		}
		defer rows.Close()

		utils.Banner("Kylrix Vault - Secrets")
		header := []string{"NAME", "CREATED"}
		var data [][]string
		for rows.Next() {
			var name, created string
			if err := rows.Scan(&name, &created); err != nil {
				return err
			}
			data = append(data, []string{name, created})
		}

		if len(data) == 0 {
			utils.Info("No secrets found in local vault.")
		} else {
			utils.Table(header, data)
		}
		return nil
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

		password, err := utils.PasswordPrompt("Vault Master Password")
		if err != nil {
			return err
		}

		salt := make([]byte, crypto.SaltSize)
		copy(salt, []byte("kylrix-ecosystem-default-salt-!!")) 

		key := crypto.DeriveKey(password, salt)
		encrypted, err := crypto.Encrypt(value, key)
		if err != nil {
			return err
		}

		database, err := db.InitDB()
		if err != nil {
			return err
		}
		defer database.Close()

		_, err = database.Exec("INSERT OR REPLACE INTO vault_secrets (name, payload) VALUES (?, ?)", name, encrypted)
		if err != nil {
			return err
		}

		utils.Success(fmt.Sprintf("Secret '%s' encrypted and saved to local SQLite vault.", name))
		return nil
	},
}

var vaultGetCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Retrieve/Decrypt a secret",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		
		database, err := db.InitDB()
		if err != nil {
			return err
		}
		defer database.Close()

		var payload string
		err = database.QueryRow("SELECT payload FROM vault_secrets WHERE name = ?", name).Scan(&payload)
		if err != nil {
			utils.Error(fmt.Sprintf("Secret '%s' not found.", name))
			return nil
		}
		
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
				utils.Error("Decryption failed.")
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
