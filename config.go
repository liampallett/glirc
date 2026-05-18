package main

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	Nick   string `json:"nick"`
	User   string `json:"user"`
	Server string `json:"server"`
	Port   int    `json:"port"`
}

func (config Config) Validate() error {
	if config.Nick == "" {
		return errors.New("config: nick is required")
	}
	if config.User == "" {
		return errors.New("config: user is required")
	}
	if config.Server == "" {
		return errors.New("config: server is required")
	}
	if config.Port == 0 {
		return errors.New("config: port is required")
	}
	return nil
}

func loadConfig(path string) (Config, error) {
	var config Config

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, config.Validate()
}
