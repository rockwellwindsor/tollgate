package state

import (
	"os"
	"testing"
)

func TestIsGlobalOn_DefaultTrue(t *testing.T) {
	m := NewManager(t.TempDir())

	on, err := m.IsGlobalOn()
	if err != nil {
		t.Fatalf("IsGlobalOn() error = %v", err)
	}
	if !on {
		t.Error("IsGlobalOn() = false, want true by default")
	}
}

func TestSessionPaused_MarkedForPID(t *testing.T) {
	m := NewManager(t.TempDir())
	pid := os.Getpid()

	if err := m.SetSessionPaused(pid); err != nil {
		t.Fatalf("SetSessionPaused() error = %v", err)
	}

	paused, err := m.IsSessionPaused(pid)
	if err != nil {
		t.Fatalf("IsSessionPaused() error = %v", err)
	}
	if !paused {
		t.Errorf("IsSessionPaused(%d) = false, want true", pid)
	}
}

func TestSetGlobalOff_IsGlobalOnReturnsFalse(t *testing.T) {
	m := NewManager(t.TempDir())

	if err := m.SetGlobalOff(); err != nil {
		t.Fatalf("SetGlobalOff() error = %v", err)
	}

	on, err := m.IsGlobalOn()
	if err != nil {
		t.Fatalf("IsGlobalOn() error = %v", err)
	}
	if on {
		t.Error("IsGlobalOn() = true after SetGlobalOff(), want false")
	}
}
