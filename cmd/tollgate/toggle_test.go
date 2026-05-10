package main

import (
	"os"
	"testing"

	"github.com/rockwellwindsor/tollgate/internal/state"
)

func TestOff_SetsGlobalOff(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TOLLGATE_HOME", dir)

	cmd := newOffCmd()
	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("off RunE error = %v", err)
	}

	mgr := state.NewManager(dir)
	on, err := mgr.IsGlobalOn()
	if err != nil {
		t.Fatalf("IsGlobalOn error = %v", err)
	}
	if on {
		t.Error("expected global to be off after tollgate off")
	}
}

func TestResume_ClearsSessionPauseMarker(t *testing.T) {
	dir := t.TempDir()
	mgr := state.NewManager(dir)
	if err := mgr.SetSessionPaused(os.Getpid()); err != nil {
		t.Fatalf("SetSessionPaused error = %v", err)
	}
	t.Setenv("TOLLGATE_HOME", dir)

	cmd := newResumeCmd()
	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("resume RunE error = %v", err)
	}

	paused, err := mgr.IsSessionPaused(os.Getpid())
	if err != nil {
		t.Fatalf("IsSessionPaused error = %v", err)
	}
	if paused {
		t.Error("expected session to not be paused after tollgate resume")
	}
}

func TestPause_SetsSessionPauseMarker(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TOLLGATE_HOME", dir)

	cmd := newPauseCmd()
	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("pause RunE error = %v", err)
	}

	mgr := state.NewManager(dir)
	paused, err := mgr.IsSessionPaused(os.Getpid())
	if err != nil {
		t.Fatalf("IsSessionPaused error = %v", err)
	}
	if !paused {
		t.Error("expected session to be paused after tollgate pause")
	}
}
