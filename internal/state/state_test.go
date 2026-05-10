package state

import (
	"os"
	"os/exec"
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

func TestDefaultDir_RespectsTollgateHome(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TOLLGATE_HOME", dir)

	if got := DefaultDir(); got != dir {
		t.Errorf("DefaultDir() = %q, want %q", got, dir)
	}
}

func TestAllowForSession_SubsequentCallSees(t *testing.T) {
	m := NewManager(t.TempDir())
	pid := os.Getpid()

	if err := m.AllowForSession("git-push", pid); err != nil {
		t.Fatalf("AllowForSession() error = %v", err)
	}

	allowed, err := m.IsAllowedForSession("git-push", pid)
	if err != nil {
		t.Fatalf("IsAllowedForSession() error = %v", err)
	}
	if !allowed {
		t.Errorf("IsAllowedForSession(%q, %d) = false, want true", "git-push", pid)
	}
}

func TestPruneStaleSessions_RemovesDeadPID(t *testing.T) {
	m := NewManager(t.TempDir())

	cmd := exec.Command("true")
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not spawn process: %v", err)
	}
	deadPID := cmd.Process.Pid

	if err := m.SetSessionPaused(deadPID); err != nil {
		t.Fatalf("SetSessionPaused() error = %v", err)
	}
	if err := m.PruneStaleSessions(); err != nil {
		t.Fatalf("PruneStaleSessions() error = %v", err)
	}

	paused, err := m.IsSessionPaused(deadPID)
	if err != nil {
		t.Fatalf("IsSessionPaused() error = %v", err)
	}
	if paused {
		t.Errorf("IsSessionPaused(%d) = true after prune, want false", deadPID)
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
