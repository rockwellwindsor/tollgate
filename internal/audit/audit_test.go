package audit

import (
	"bufio"
	"encoding/json"
	"fmt"
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

func TestWrite_ValidJSONLinesPerEntry(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	entries := []Entry{
		{Binary: "git", Args: []string{"push"}, Pattern: "git-push", Decision: "allowed-once"},
		{Binary: "gh", Args: []string{"pr", "create"}, Pattern: "gh-pr-create", Decision: "denied"},
	}
	for _, e := range entries {
		if err := Write(path, e); err != nil {
			t.Fatalf("Write() error = %v", err)
		}
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	lineCount := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if !json.Valid(line) {
			t.Errorf("line %d is not valid JSON: %s", lineCount+1, line)
		}
		lineCount++
	}
	if lineCount != len(entries) {
		t.Errorf("got %d lines, want %d", lineCount, len(entries))
	}
}

func TestWrite_CreatesDirectoryIfMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "nested", "audit.log")

	entry := Entry{Binary: "git", Args: []string{"push"}, Pattern: "git-push", Decision: "allowed-once"}

	if err := Write(path, entry); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("log file not created: %v", err)
	}
}

func TestWrite_ConcurrentWritesProduceValidLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	const n = 20
	errc := make(chan error, n)
	for i := range n {
		go func(i int) {
			errc <- Write(path, Entry{Binary: "git", Args: []string{"push"}, Pattern: "git-push", Decision: fmt.Sprintf("decision-%d", i)})
		}(i)
	}
	for range n {
		if err := <-errc; err != nil {
			t.Errorf("Write() error = %v", err)
		}
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	lineCount := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if !json.Valid(line) {
			t.Errorf("line %d is not valid JSON: %s", lineCount+1, line)
		}
		lineCount++
	}
	if lineCount != n {
		t.Errorf("got %d lines, want %d", lineCount, n)
	}
}
