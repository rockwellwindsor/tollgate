package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultAction string `json:"default_action"`
}

func DefaultPath() string {
	if home := os.Getenv("TOLLGATE_HOME"); home != "" {
		return filepath.Join(home, "config.json")
	}
	return filepath.Join(os.Getenv("HOME"), ".tollgate", "config.json")
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
