package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rockwellwindsor/tollgate/internal/testkit"
)

var shimGitPath string

func TestMain(m *testing.M) {
	bin, err := buildShim("./cmd/shim-git")
	if err != nil {
		panic("failed to build shim-git: " + err.Error())
	}
	shimGitPath = bin
	os.Exit(m.Run())
}

func buildShim(pkg string) (string, error) {
	dir, err := os.MkdirTemp("", "tollgate-e2e-*")
	if err != nil {
		return "", err
	}
	out := filepath.Join(dir, "shim-git")
	cmd := exec.Command("go", "build", "-o", out, pkg)
	cmd.Dir = filepath.Join(os.Getenv("PWD"), "../..")
	if b, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build failed: %s", b)
	}
	return out, nil
}

func TestPassthrough(t *testing.T) {
	home := testkit.TempHome(t)

	realDir := t.TempDir()
	argsFile := filepath.Join(home, "git-args")
	fakeGit := testkit.FakeBinary(t, realDir, "git")
	_ = fakeGit

	script := "#!/bin/sh\necho \"$@\" > " + argsFile + "\n"
	if err := os.WriteFile(filepath.Join(realDir, "git"), []byte(script), 0755); err != nil {
		t.Fatal(err)
	}

	shimDir := filepath.Dir(shimGitPath)
	cmd := exec.Command(shimGitPath, "push", "origin", "main")
	cmd.Env = append(os.Environ(),
		"PATH="+shimDir+string(os.PathListSeparator)+realDir,
		"HOME="+home,
	)

	if err := cmd.Run(); err != nil {
		t.Fatalf("shim exited with error: %v", err)
	}

	got, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatalf("fake git was never called — args file not found: %v", err)
	}

	want := "push origin main"
	if gotStr := strings.TrimSpace(string(got)); gotStr != want {
		t.Errorf("fake git called with %q, want %q", gotStr, want)
	}
}

func TestExitCodePassthrough(t *testing.T) {
	realDir := t.TempDir()
	script := "#!/bin/sh\nexit 42\n"
	if err := os.WriteFile(filepath.Join(realDir, "git"), []byte(script), 0755); err != nil {
		t.Fatal(err)
	}

	shimDir := filepath.Dir(shimGitPath)
	cmd := exec.Command(shimGitPath, "push")
	cmd.Env = append(os.Environ(),
		"PATH="+shimDir+string(os.PathListSeparator)+realDir,
	)

	err := cmd.Run()
	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	}

	if exitCode != 42 {
		t.Errorf("exit code = %d, want 42", exitCode)
	}
}
