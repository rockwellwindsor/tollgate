package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rockwellwindsor/tollgate/internal/audit"
)

func writeAuditLog(t *testing.T, dir string) string {
	t.Helper()
	logPath := filepath.Join(dir, "audit.log")
	e := audit.Entry{Binary: "git", Args: []string{"push"}, Pattern: "git-push", Decision: "allowed-once"}
	if err := audit.Write(logPath, e); err != nil {
		t.Fatalf("audit.Write() error = %v", err)
	}
	return logPath
}

func TestClearLogs_DoesNotTouchOtherFiles(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TOLLGATE_HOME", dir)
	writeAuditLog(t, dir)

	configPath := filepath.Join(dir, "config.json")
	if err := os.WriteFile(configPath, []byte(`{}`), 0644); err != nil {
		t.Fatalf("WriteFile config error = %v", err)
	}
	shimPath := filepath.Join(dir, "shim-git")
	if err := os.WriteFile(shimPath, []byte(`#!/bin/sh`), 0755); err != nil {
		t.Fatalf("WriteFile shim error = %v", err)
	}

	cmd := newClearLogsCmd()
	cmd.SetOut(&bytes.Buffer{})
	if err := cmd.ParseFlags([]string{"--yes"}); err != nil {
		t.Fatalf("ParseFlags error = %v", err)
	}
	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("RunE error = %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("clear-logs deleted config.json — it should not have")
	}
	if _, err := os.Stat(shimPath); os.IsNotExist(err) {
		t.Error("clear-logs deleted shim-git — it should not have")
	}
}

func TestClearLogs_PromptNoKeepsLog(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TOLLGATE_HOME", dir)
	logPath := writeAuditLog(t, dir)

	cmd := newClearLogsCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetIn(strings.NewReader("n\n"))
	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("RunE error = %v", err)
	}

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("audit.log should not be deleted after answering n")
	}
}

func TestClearLogs_PromptYesDeletesLog(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TOLLGATE_HOME", dir)
	logPath := writeAuditLog(t, dir)

	cmd := newClearLogsCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetIn(strings.NewReader("y\n"))
	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("RunE error = %v", err)
	}

	if _, err := os.Stat(logPath); !os.IsNotExist(err) {
		t.Error("expected audit.log to be deleted after answering y")
	}
}

func TestClearLogs_YesDeletesLog(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TOLLGATE_HOME", dir)
	logPath := writeAuditLog(t, dir)

	cmd := newClearLogsCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	if err := cmd.ParseFlags([]string{"--yes"}); err != nil {
		t.Fatalf("ParseFlags error = %v", err)
	}
	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("RunE error = %v", err)
	}

	if _, err := os.Stat(logPath); !os.IsNotExist(err) {
		t.Error("expected audit.log to be deleted after --yes")
	}
}

func TestClearLogs_DryRun(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TOLLGATE_HOME", dir)
	logPath := writeAuditLog(t, dir)

	cmd := newClearLogsCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	if err := cmd.ParseFlags([]string{"--dry-run"}); err != nil {
		t.Fatalf("ParseFlags error = %v", err)
	}
	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("RunE error = %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "audit.log") {
		t.Errorf("dry-run output %q does not mention audit.log", out)
	}
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("dry-run deleted audit.log — it should not have")
	}
}
