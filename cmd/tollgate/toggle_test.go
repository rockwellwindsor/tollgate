package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/rockwellwindsor/tollgate/internal/state"
)

func TestOn_SetsGlobalOn(t *testing.T) {
	dir := t.TempDir()
	mgr := state.NewManager(dir)
	if err := mgr.SetGlobalOff(); err != nil {
		t.Fatalf("SetGlobalOff error = %v", err)
	}
	t.Setenv("TOLLGATE_HOME", dir)

	cmd := newOnCmd()
	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("on RunE error = %v", err)
	}

	on, err := mgr.IsGlobalOn()
	if err != nil {
		t.Fatalf("IsGlobalOn error = %v", err)
	}
	if !on {
		t.Error("expected global to be on after tollgate on")
	}
}

func TestOn_WhenAlreadyOn_PrintsAlreadyEnabled(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TOLLGATE_HOME", dir)

	cmd := newOnCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("tollgate on (already on): error = %v", err)
	}
	if !strings.Contains(buf.String(), "already enabled") {
		t.Errorf("got %q, want output containing \"already enabled\"", buf.String())
	}
}

func TestOff_WhenAlreadyOff_PrintsAlreadyDisabled(t *testing.T) {
	dir := t.TempDir()
	mgr := state.NewManager(dir)
	if err := mgr.SetGlobalOff(); err != nil {
		t.Fatalf("SetGlobalOff error = %v", err)
	}
	t.Setenv("TOLLGATE_HOME", dir)

	cmd := newOffCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("tollgate off (already off): error = %v", err)
	}
	if !strings.Contains(buf.String(), "already disabled") {
		t.Errorf("got %q, want output containing \"already disabled\"", buf.String())
	}
}

func TestPause_WhenAlreadyPaused_PrintsAlreadyPaused(t *testing.T) {
	dir := t.TempDir()
	mgr := state.NewManager(dir)
	if err := mgr.SetSessionPaused(os.Getpid()); err != nil {
		t.Fatalf("SetSessionPaused error = %v", err)
	}
	t.Setenv("TOLLGATE_HOME", dir)

	cmd := newPauseCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("tollgate pause (already paused): error = %v", err)
	}
	if !strings.Contains(buf.String(), "already paused") {
		t.Errorf("got %q, want output containing \"already paused\"", buf.String())
	}
}

func TestResume_WhenNotPaused_PrintsNotPaused(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TOLLGATE_HOME", dir)

	cmd := newResumeCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("tollgate resume (not paused): error = %v", err)
	}
	if !strings.Contains(buf.String(), "not paused") {
		t.Errorf("got %q, want output containing \"not paused\"", buf.String())
	}
}

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
