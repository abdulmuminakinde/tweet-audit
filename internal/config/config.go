package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey string
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

func (c *Config) SetAPIKey(apikey string) error {
	c.APIKey = apikey
	dat, err := json.Marshal(c)
	if err != nil {
		return err
	}

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(path, dat, 0644)
	if err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	conf := Config{}
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	dat, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(dat, &conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
