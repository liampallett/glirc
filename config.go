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

	if config.Nick == "" {
		return Config{}, errors.New("config: nick is required")
	}
	if config.User == "" {
		return Config{}, errors.New("config: user is required")
	}
	if config.Server == "" {
		return Config{}, errors.New("config: server is required")
	}
	if config.Port == 0 {
		return Config{}, errors.New("config: port is required")
	}

	return config, nil
}
