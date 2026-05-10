package testkit

import (
	"os"
	"path/filepath"
	"testing"
)

// FakeBinary writes an executable shell script named `name` into `dir`.
// body is the script content after the shebang line.
// Returns the full path to the created file.
func FakeBinary(t *testing.T, dir, name, body string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	script := "#!/bin/sh\n" + body + "\n"
	if err := os.WriteFile(path, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	return path
}
