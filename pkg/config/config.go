package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	BaseURI string `json:"base_uri"`
	APIKey  string `json:"api_key"`
}

func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".kylrix"), nil
}
