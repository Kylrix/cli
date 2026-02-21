package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/nathfavour/kylrix/cli/pkg/config"
	"github.com/nathfavour/kylrix/cli/pkg/crypto"
	"github.com/nathfavour/kylrix/cli/pkg/db"
	"github.com/nathfavour/kylrix/cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	decryptSecret bool
)

// getMEK handles the multi-layered unlocking logic: Ephemeral PIN -> Master Password
func getMEK(cfg *config.Config) ([]byte, error) {
	// 1. Try Ephemeral PIN first if available
	if cfg.EphemeralSession != nil && cfg.PinVerifier != nil {
		pin, err := utils.PasswordPrompt("Enter 4-digit PIN to unlock")
		if err == nil && len(pin) == 4 {
			// Verify PIN hash
			salt, _ := base64.StdEncoding.DecodeString(cfg.PinVerifier.Salt)
			expectedHash, _ := base64.StdEncoding.DecodeString(cfg.PinVerifier.Hash)
			actualHash := crypto.DerivePinKey(pin, salt)

			if string(actualHash) == string(expectedHash) {
				// PIN correct, unwrap MEK
				sessionSalt, _ := base64.StdEncoding.DecodeString(cfg.EphemeralSession.SessionSalt)
				ephemeralKey := crypto.DeriveEphemeralKey(pin, sessionSalt)
				mek, err := crypto.UnwrapKey(cfg.EphemeralSession.WrappedMek, ephemeralKey)
				if err == nil {
					utils.Success("Vault unlocked via Ephemeral PIN.")
					return mek, nil
				}
			}
			utils.Warning("PIN incorrect or session expired.")
		}
	}

	// 2. Fallback to Master Password
	password, err := utils.PasswordPrompt("Vault Master Password")
	if err != nil {
		return nil, err
	}

	salt := make([]byte, crypto.SaltSize)
	copy(salt, []byte("kylrix-ecosystem-default-salt-!!"))
	mek := crypto.DeriveKey(password, salt)

	// 3. If PIN is set, piggyback this session
	if cfg.PinVerifier != nil {
		pin, err := utils.PasswordPrompt("Enter 4-digit PIN to secure this session")
		if err == nil && len(pin) == 4 {
			sessionSalt := make([]byte, crypto.SessionSaltSize)
			rand.Read(sessionSalt)

			ephemeralKey := crypto.DeriveEphemeralKey(pin, sessionSalt)
			wrappedMek, err := crypto.WrapKey(mek, ephemeralKey)
			if err == nil {
				cfg.EphemeralSession = &config.EphemeralSession{
					WrappedMek:  wrappedMek,
					SessionSalt: base64.StdEncoding.EncodeToString(sessionSalt),
				}
				config.SaveConfig(cfg)
				utils.Success("Session piggybacked with PIN.")
			}
		}
	}

	return mek, nil
}

var vaultSetupPinCmd = &cobra.Command{
	Use:   "setup-pin",
	Short: "Setup a 4-digit PIN for quick unlocking",
	RunE: func(cmd *cobra.Command, args []string) error {
		utils.Banner("Kylrix Vault - Setup PIN")
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		pin, err := utils.PasswordPrompt("Choose a 4-digit PIN")
		if err != nil || len(pin) != 4 {
			return fmt.Errorf("invalid PIN: must be 4 digits")
		}

		salt := make([]byte, crypto.PinSaltSize)
		rand.Read(salt)
		hash := crypto.DerivePinKey(pin, salt)

		cfg.PinVerifier = &config.PinVerifier{
			Salt: base64.StdEncoding.EncodeToString(salt),
			Hash: base64.StdEncoding.EncodeToString(hash),
		}

		err = config.SaveConfig(cfg)
		if err != nil {
			return err
		}

		utils.Success("PIN verifier setup on disk. Next login will allow piggybacking.")
		return nil
	},
}

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

		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		key, err := getMEK(cfg)
		if err != nil {
			return err
		}

		encrypted, err := crypto.Encrypt(value, key)
		if err != nil {
			return err
		}
		// Explicitly zero the MEK after use
		crypto.ZeroBytes(key)

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
			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}

			key, err := getMEK(cfg)
			if err != nil {
				return err
			}

			decrypted, err := crypto.Decrypt(payload, key)
			if err != nil {
				utils.Error("Decryption failed.")
				return err
			}
			// Explicitly zero the MEK after use
			crypto.ZeroBytes(key)

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
	vaultCmd.AddCommand(vaultSetupPinCmd)
	rootCmd.AddCommand(vaultCmd)
}
