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
