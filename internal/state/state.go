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

func (m *Manager) sessionAllowPath(pattern string, pid int) string {
	return filepath.Join(m.dir, fmt.Sprintf("allowed-%d-%s", pid, pattern))
}

func (m *Manager) AllowForSession(pattern string, pid int) error {
	if err := os.MkdirAll(m.dir, 0755); err != nil {
		return err
	}
	f, err := os.Create(m.sessionAllowPath(pattern, pid))
	if err != nil {
		return err
	}
	return f.Close()
}

func (m *Manager) IsAllowedForSession(pattern string, pid int) (bool, error) {
	_, err := os.Stat(m.sessionAllowPath(pattern, pid))
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
