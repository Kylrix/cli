package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	BaseURI string `json:"base_uri"`
	APIKey  string `json:"api_key"`
	Token   string `json:"token"`
}

func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".kylrix")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return "", err
		}
	}
	return dir, nil
}

func GetConfigFile() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
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
