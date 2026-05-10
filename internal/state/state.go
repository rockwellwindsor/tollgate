package state

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type Manager struct {
	dir string
}

func NewManager(dir string) *Manager {
	return &Manager{dir: dir}
}

func DefaultDir() string {
	if home := os.Getenv("TOLLGATE_HOME"); home != "" {
		return home
	}
	return filepath.Join(os.Getenv("HOME"), ".tollgate")
}

func (m *Manager) createSentinel(name string) error {
	if err := os.MkdirAll(m.dir, 0755); err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(m.dir, name))
	if err != nil {
		return err
	}
	return f.Close()
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
	return m.createSentinel("disabled")
}

func (m *Manager) SetSessionPaused(pid int) error {
	return m.createSentinel(fmt.Sprintf("paused-%d", pid))
}

func (m *Manager) IsSessionPaused(pid int) (bool, error) {
	_, err := os.Stat(filepath.Join(m.dir, fmt.Sprintf("paused-%d", pid)))
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return err == nil, err
}

func (m *Manager) AllowForSession(pattern string, pid int) error {
	return m.createSentinel(fmt.Sprintf("allowed-%d-%s", pid, pattern))
}

func (m *Manager) IsAllowedForSession(pattern string, pid int) (bool, error) {
	_, err := os.Stat(filepath.Join(m.dir, fmt.Sprintf("allowed-%d-%s", pid, pattern)))
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return err == nil, err
}

func (m *Manager) PruneStaleSessions() error {
	entries, err := os.ReadDir(m.dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasPrefix(name, "paused-") {
			continue
		}
		pid, err := strconv.Atoi(strings.TrimPrefix(name, "paused-"))
		if err != nil {
			continue
		}
		if !pidAlive(pid) {
			os.Remove(filepath.Join(m.dir, name))
		}
	}
	return nil
}

func pidAlive(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return p.Signal(syscall.Signal(0)) == nil
}
