package state

import (
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
