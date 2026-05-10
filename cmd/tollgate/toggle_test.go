package main

import (
	"os"
	"testing"

	"github.com/rockwellwindsor/tollgate/internal/state"
)

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
