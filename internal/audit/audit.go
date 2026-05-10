package audit

import (
	"encoding/json"
	"errors"
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
	return nil, errors.New("not implemented")
}
