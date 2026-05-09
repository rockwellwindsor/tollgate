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
	if !strings.Contains(got, "0.1.0") {
		t.Errorf("--version output %q does not contain version string", got)
	}
}
