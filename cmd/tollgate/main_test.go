package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	cmd := newRootCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--version"})
	_ = cmd.Execute()

	got := buf.String()
	if !strings.Contains(got, "tollgate version") {
		t.Errorf("--version output %q does not contain version string", got)
	}
}

func TestHelpListsSubcommands(t *testing.T) {
	cmd := newRootCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})
	_ = cmd.Execute()

	got := buf.String()
	for _, want := range []string{"tollgate", "status", "pause", "install"} {
		if !strings.Contains(got, want) {
			t.Errorf("--help output missing %q", want)
		}
	}
}
