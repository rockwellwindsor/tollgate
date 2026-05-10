package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	DefaultAction string `json:"default_action"`
}

func defaults() Config {
	return Config{DefaultAction: "prompt"}
}

func Load(path string) (Config, error) {
	cfg := defaults()

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return Config{}, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
