package config

import (
	"path/filepath"
	"testing"
)

func TestLoad_MissingFile_ReturnsDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.DefaultAction != "prompt" {
		t.Errorf("DefaultAction = %q, want %q", cfg.DefaultAction, "prompt")
	}
}
