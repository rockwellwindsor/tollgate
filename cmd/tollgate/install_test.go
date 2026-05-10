package main

import (
	"bytes"
	"os"
	"path/filepath"
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
