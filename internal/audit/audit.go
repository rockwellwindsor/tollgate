package audit

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

type Entry struct {
	Binary   string   `json:"binary"`
	Args     []string `json:"args"`
	Pattern  string   `json:"pattern"`
	Decision string   `json:"decision"`
}

func Write(path string, entry Entry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	line, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = f.Write(append(line, '\n'))
	return err
}

func Read(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var e Entry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, scanner.Err()
}
