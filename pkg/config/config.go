package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	BaseURI          string            `json:"base_uri"`
	APIKey           string            `json:"api_key"`
	Token            string            `json:"token"`
	PinVerifier      *PinVerifier      `json:"pin_verifier,omitempty"`
	EphemeralSession *EphemeralSession `json:"ephemeral_session,omitempty"`
}

type PinVerifier struct {
	Salt string `json:"salt"`
	Hash string `json:"hash"`
}

type EphemeralSession struct {
	WrappedMek  string `json:"wrapped_mek"`
	SessionSalt string `json:"session_salt"`
}

func GetAppConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(configDir, "kylrix")
	
	// Create subdirectories
	subdirs := []string{"configs", "data", "logs", "cache"}
	for _, d := range subdirs {
		path := filepath.Join(appDir, d)
		if err := os.MkdirAll(path, 0755); err != nil {
			return "", err
		}
	}
	
	return appDir, nil
}

func GetConfigFile() (string, error) {
	appDir, err := GetAppConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appDir, "configs", "config.json"), nil
}

func GetDataDir() (string, error) {
	appDir, err := GetAppConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appDir, "data"), nil
}

func LoadConfig() (*Config, error) {
	file, err := GetConfigFile()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return &Config{BaseURI: "https://api.kylrix.com"}, nil
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	file, err := GetConfigFile()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0644)
}
