package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rockwellwindsor/tollgate/internal/audit"
	"github.com/rockwellwindsor/tollgate/internal/state"
)

func runStatus(t *testing.T) string {
	t.Helper()
	cmd := newStatusCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	_ = cmd.RunE(cmd, []string{})
	return buf.String()
}

func TestStatus_SessionAllowPatterns(t *testing.T) {
	dir := t.TempDir()
	mgr := state.NewManager(dir)
	if err := mgr.AllowForSession("git-push", os.Getpid()); err != nil {
		t.Fatalf("AllowForSession() error = %v", err)
	}
	t.Setenv("TOLLGATE_HOME", dir)

	out := runStatus(t)
	if !strings.Contains(out, "git-push") {
		t.Errorf("status output %q does not contain allowed pattern", out)
	}
}

func TestStatus_AuditEntryCount(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TOLLGATE_HOME", dir)

	logPath := filepath.Join(dir, "audit.log")
	entries := []audit.Entry{
		{Binary: "git", Args: []string{"push"}, Pattern: "git-push", Decision: "allowed-once"},
		{Binary: "gh", Args: []string{"pr", "create"}, Pattern: "gh-pr-create", Decision: "denied"},
	}
	for _, e := range entries {
		if err := audit.Write(logPath, e); err != nil {
			t.Fatalf("audit.Write() error = %v", err)
		}
	}

	out := runStatus(t)
	if !strings.Contains(out, "2") {
		t.Errorf("status output %q does not contain audit entry count", out)
	}
}

func TestStatus_SessionPaused(t *testing.T) {
	dir := t.TempDir()
	mgr := state.NewManager(dir)
	if err := mgr.SetSessionPaused(os.Getpid()); err != nil {
		t.Fatalf("SetSessionPaused() error = %v", err)
	}
	t.Setenv("TOLLGATE_HOME", dir)

	out := runStatus(t)
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

	out := runStatus(t)
	if !strings.Contains(out, "DISABLED") {
		t.Errorf("status output %q does not contain DISABLED", out)
	}
}

func TestStatus_GlobalEnabled(t *testing.T) {
	dir := t.TempDir()
	_ = state.NewManager(dir) // no SetGlobalOff — defaults to on
	t.Setenv("TOLLGATE_HOME", dir)

	out := runStatus(t)
	if !strings.Contains(out, "ENABLED") {
		t.Errorf("status output %q does not contain ENABLED", out)
	}
}
