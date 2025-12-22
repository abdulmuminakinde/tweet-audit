package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey   string
	Username string
}

const configFile = ".tweetauditconfig.json"

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, configFile)
	return path, nil
}

func LoadOrCreateConfig() (*Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	dat, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	var conf Config
	err = json.Unmarshal(dat, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func (c *Config) Save() error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	dat, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, dat, 0644)
}
