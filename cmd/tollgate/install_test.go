package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fakeShimSrcDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, name := range []string{"shim-git", "shim-gh"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("#!/bin/sh\n"), 0755); err != nil {
			t.Fatalf("WriteFile %s error = %v", name, err)
		}
	}
	return dir
}

func TestInstall_Idempotent(t *testing.T) {
	home := t.TempDir()
	t.Setenv("TOLLGATE_HOME", home)
	srcDir := fakeShimSrcDir(t)

	cmd := newInstallCmd()
	cmd.SetOut(&bytes.Buffer{})
	if err := install(cmd, srcDir); err != nil {
		t.Fatalf("first install error = %v", err)
	}
	if err := install(cmd, srcDir); err != nil {
		t.Fatalf("second install error = %v", err)
	}

	for _, name := range []string{"shim-git", "shim-gh"} {
		if _, err := os.Stat(filepath.Join(home, "bin", name)); os.IsNotExist(err) {
			t.Errorf("%s missing after second install", name)
		}
	}
}

func TestInstall_PrintsPathInstructions(t *testing.T) {
	home := t.TempDir()
	t.Setenv("TOLLGATE_HOME", home)
	srcDir := fakeShimSrcDir(t)

	cmd := newInstallCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	if err := install(cmd, srcDir); err != nil {
		t.Fatalf("install error = %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "PATH") {
		t.Errorf("install output %q does not mention PATH", out)
	}
	if !strings.Contains(out, filepath.Join(home, "bin")) {
		t.Errorf("install output %q does not contain bin dir path", out)
	}
}

func TestInstall_LeavesExistingConfigAlone(t *testing.T) {
	home := t.TempDir()
	t.Setenv("TOLLGATE_HOME", home)
	srcDir := fakeShimSrcDir(t)

	configPath := filepath.Join(home, "config.json")
	existing := []byte(`{"default_action":"deny"}`)
	if err := os.WriteFile(configPath, existing, 0644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	cmd := newInstallCmd()
	cmd.SetOut(&bytes.Buffer{})
	if err := install(cmd, srcDir); err != nil {
		t.Fatalf("install error = %v", err)
	}

	got, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}
	if string(got) != string(existing) {
		t.Errorf("config.json was modified: got %q, want %q", got, existing)
	}
}

func TestInstall_WritesDefaultConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("TOLLGATE_HOME", home)
	srcDir := fakeShimSrcDir(t)

	cmd := newInstallCmd()
	cmd.SetOut(&bytes.Buffer{})
	if err := install(cmd, srcDir); err != nil {
		t.Fatalf("install error = %v", err)
	}

	configPath := filepath.Join(home, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("expected config.json to be written by install")
	}
}

func TestInstall_WritesShimsWithExecuteBits(t *testing.T) {
	home := t.TempDir()
	t.Setenv("TOLLGATE_HOME", home)
	srcDir := fakeShimSrcDir(t)

	cmd := newInstallCmd()
	cmd.SetOut(&bytes.Buffer{})
	if err := install(cmd, srcDir); err != nil {
		t.Fatalf("install error = %v", err)
	}

	for _, name := range []string{"shim-git", "shim-gh"} {
		path := filepath.Join(home, "bin", name)
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			t.Errorf("expected %s to exist", path)
			continue
		}
		if info.Mode()&0111 == 0 {
			t.Errorf("%s is not executable (mode %o)", name, info.Mode())
		}
	}
}

func TestInstall_CreatesBinDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("TOLLGATE_HOME", home)
	srcDir := fakeShimSrcDir(t)

	cmd := newInstallCmd()
	cmd.SetOut(&bytes.Buffer{})
	if err := install(cmd, srcDir); err != nil {
		t.Fatalf("install error = %v", err)
	}

	binDir := filepath.Join(home, "bin")
	if _, err := os.Stat(binDir); os.IsNotExist(err) {
		t.Errorf("expected %s to exist after install", binDir)
	}
}
