package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/rockwellwindsor/tollgate/internal/state"
)

func runStatus(t *testing.T, stateDir, auditPath string) string {
	t.Helper()
	cmd := newStatusCmd()
	cmd.SetArgs([]string{})
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	_ = cmd.RunE(cmd, []string{})
	return buf.String()
}

func TestStatus_SessionPaused(t *testing.T) {
	dir := t.TempDir()
	mgr := state.NewManager(dir)
	if err := mgr.SetSessionPaused(os.Getpid()); err != nil {
		t.Fatalf("SetSessionPaused() error = %v", err)
	}
	t.Setenv("TOLLGATE_HOME", dir)

	out := runStatus(t, dir, "")
	if !strings.Contains(out, "paused") {
		t.Errorf("status output %q does not contain paused", out)
	}
}

func TestStatus_GlobalDisabled(t *testing.T) {
	dir := t.TempDir()
	mgr := state.NewManager(dir)
	if err := mgr.SetGlobalOff(); err != nil {
		t.Fatalf("SetGlobalOff() error = %v", err)
	}
	t.Setenv("TOLLGATE_HOME", dir)

	out := runStatus(t, dir, "")
	if !strings.Contains(out, "DISABLED") {
		t.Errorf("status output %q does not contain DISABLED", out)
	}
}

func TestStatus_GlobalEnabled(t *testing.T) {
	dir := t.TempDir()
	_ = state.NewManager(dir) // no SetGlobalOff — defaults to on
	t.Setenv("TOLLGATE_HOME", dir)

	out := runStatus(t, dir, "")
	if !strings.Contains(out, "ENABLED") {
		t.Errorf("status output %q does not contain ENABLED", out)
	}
}
