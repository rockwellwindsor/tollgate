package state

import (
	"errors"
	"os"
	"path/filepath"
)

type Manager struct {
	dir string
}

func NewManager(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) IsGlobalOn() (bool, error) {
	_, err := os.Stat(filepath.Join(m.dir, "disabled"))
	if errors.Is(err, os.ErrNotExist) {
		return true, nil
	}
	return false, err
}

func (m *Manager) SetGlobalOn() error {
	return os.Remove(filepath.Join(m.dir, "disabled"))
}

func (m *Manager) SetGlobalOff() error {
	if err := os.MkdirAll(m.dir, 0755); err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(m.dir, "disabled"))
	if err != nil {
		return err
	}
	return f.Close()
}
