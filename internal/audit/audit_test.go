package audit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWrite_AppendsJSONLine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	entry := Entry{Binary: "git", Args: []string{"push", "origin", "main"}, Pattern: "git-push", Decision: "allowed-once"}

	if err := Write(path, entry); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read log file: %v", err)
	}
	if len(data) == 0 {
		t.Error("log file is empty, expected a JSON line")
	}
}
