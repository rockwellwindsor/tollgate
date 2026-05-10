package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rockwellwindsor/tollgate/internal/state"
	"github.com/rockwellwindsor/tollgate/internal/testkit"
)

func newShimCmd(t *testing.T, home, realDir string, args ...string) *exec.Cmd {
	t.Helper()
	return newShimCmdFor(t, shimGitPath, home, realDir, args...)
}

func TestPrompt_AllowSessionInvokesFakeGitAndRecordsAllow(t *testing.T) {
	home := testkit.TempHome(t)
	realDir := t.TempDir()
	argsFile := filepath.Join(home, "git-args")
	testkit.FakeBinary(t, realDir, "git", `echo "$@" > `+argsFile)

	cmd := newShimCmd(t, home, realDir, "push", "origin", "main")
	cmd.Stdin = strings.NewReader("a")

	if err := cmd.Run(); err != nil {
		t.Fatalf("shim exited with error: %v", err)
	}

	got, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatalf("fake git was never called — args file not found: %v", err)
	}
	if want := "push origin main"; strings.TrimSpace(string(got)) != want {
		t.Errorf("fake git args = %q, want %q", strings.TrimSpace(string(got)), want)
	}

	mgr := state.NewManager(home)
	allowed, err := mgr.IsAllowedForSession("git-push", os.Getpid())
	if err != nil {
		t.Fatalf("IsAllowedForSession() error = %v", err)
	}
	if !allowed {
		t.Error("session allow not recorded after 'a' response")
	}
}

func TestPrompt_SessionAllowSkipsPrompt(t *testing.T) {
	home := testkit.TempHome(t)
	realDir := t.TempDir()
	argsFile := filepath.Join(home, "git-args")
	testkit.FakeBinary(t, realDir, "git", `echo "$@" > `+argsFile)

	mgr := state.NewManager(home)
	if err := mgr.AllowForSession("git-push", os.Getpid()); err != nil {
		t.Fatalf("AllowForSession() error = %v", err)
	}

	cmd := newShimCmd(t, home, realDir, "push", "origin", "main")
	// no stdin — if shim tries to prompt it will error and exit non-zero

	if err := cmd.Run(); err != nil {
		t.Fatalf("shim exited with error: %v", err)
	}

	got, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatalf("fake git was never called — args file not found: %v", err)
	}
	if want := "push origin main"; strings.TrimSpace(string(got)) != want {
		t.Errorf("fake git args = %q, want %q", strings.TrimSpace(string(got)), want)
	}
}

func TestPrompt_TollgateOffBypassesEverything(t *testing.T) {
	home := testkit.TempHome(t)
	realDir := t.TempDir()
	argsFile := filepath.Join(home, "git-args")
	testkit.FakeBinary(t, realDir, "git", `echo "$@" > `+argsFile)

	cmd := newShimCmd(t, home, realDir, "push", "origin", "main")
	cmd.Env = append(cmd.Env, "TOLLGATE=off")
	// no stdin — bypass must happen before any prompt attempt

	if err := cmd.Run(); err != nil {
		t.Fatalf("shim exited with error: %v", err)
	}

	got, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatalf("fake git was never called — args file not found: %v", err)
	}
	if want := "push origin main"; strings.TrimSpace(string(got)) != want {
		t.Errorf("fake git args = %q, want %q", strings.TrimSpace(string(got)), want)
	}
}

func TestPrompt_NoDeniesAndExitsNonZero(t *testing.T) {
	home := testkit.TempHome(t)
	realDir := t.TempDir()
	argsFile := filepath.Join(home, "git-args")
	testkit.FakeBinary(t, realDir, "git", `echo "$@" > `+argsFile)

	cmd := newShimCmd(t, home, realDir, "push", "origin", "main")
	cmd.Stdin = strings.NewReader("n")

	err := cmd.Run()
	if err == nil {
		t.Fatal("shim exited zero after 'n', want non-zero")
	}

	if _, statErr := os.Stat(argsFile); statErr == nil {
		t.Error("fake git was called after 'n', want no invocation")
	}
}

func TestPrompt_YesInvokesFakeGit(t *testing.T) {
	home := testkit.TempHome(t)
	realDir := t.TempDir()
	argsFile := filepath.Join(home, "git-args")
	testkit.FakeBinary(t, realDir, "git", `echo "$@" > `+argsFile)

	cmd := newShimCmd(t, home, realDir, "push", "origin", "main")
	cmd.Stdin = strings.NewReader("y")

	if err := cmd.Run(); err != nil {
		t.Fatalf("shim exited with error: %v", err)
	}

	got, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatalf("fake git was never called — args file not found: %v", err)
	}
	if want := "push origin main"; strings.TrimSpace(string(got)) != want {
		t.Errorf("fake git args = %q, want %q", strings.TrimSpace(string(got)), want)
	}
}
