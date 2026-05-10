package config

import (
	"os"
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

func TestDefaultPath_RespectsTollgateHome(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TOLLGATE_HOME", dir)

	got := DefaultPath()
	want := filepath.Join(dir, "config.json")
	if got != want {
		t.Errorf("DefaultPath() = %q, want %q", got, want)
	}
}

func TestLoad_MalformedJSON_ReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{not valid json`), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() expected error for malformed JSON, got nil")
	}
}

func TestLoad_UnknownFields_NoError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"default_action":"allow","future_feature":true}`), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if cfg.DefaultAction != "allow" {
		t.Errorf("DefaultAction = %q, want %q", cfg.DefaultAction, "allow")
	}
}

func TestLoad_ValidFile_ReturnsParsedValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"default_action":"deny"}`), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.DefaultAction != "deny" {
		t.Errorf("DefaultAction = %q, want %q", cfg.DefaultAction, "deny")
	}
}
