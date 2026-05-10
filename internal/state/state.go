package state

import (
	"errors"
	"fmt"
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

func (m *Manager) sessionPausedPath(pid int) string {
	return filepath.Join(m.dir, fmt.Sprintf("paused-%d", pid))
}

func (m *Manager) SetSessionPaused(pid int) error {
	if err := os.MkdirAll(m.dir, 0755); err != nil {
		return err
	}
	f, err := os.Create(m.sessionPausedPath(pid))
	if err != nil {
		return err
	}
	return f.Close()
}

func (m *Manager) IsSessionPaused(pid int) (bool, error) {
	_, err := os.Stat(m.sessionPausedPath(pid))
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return err == nil, err
}
