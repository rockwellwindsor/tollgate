package testkit

import (
	"os"
	"path/filepath"
	"testing"
)

// FakeBinary writes an executable shell script named `name` into `dir`.
// Returns the full path to the created file.
func FakeBinary(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatal(err)
	}
	return path
}
