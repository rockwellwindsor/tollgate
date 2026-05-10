package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rockwellwindsor/tollgate/internal/testkit"
)

var shimGhPath string

func init() {
	bin, err := buildShim("./cmd/shim-gh")
	if err != nil {
		panic("failed to build shim-gh: " + err.Error())
	}
	shimGhPath = bin
}

func newShimGhCmd(t *testing.T, home, realDir string, args ...string) *exec.Cmd {
	t.Helper()
	shimDir := filepath.Dir(shimGhPath)
	cmd := exec.Command(shimGhPath, args...)
	cmd.Env = append(os.Environ(),
		"PATH="+shimDir+string(os.PathListSeparator)+realDir,
		"HOME="+home,
		"TOLLGATE_HOME="+home,
	)
	return cmd
}

func TestGhShim_PrList_PassesThrough(t *testing.T) {
	home := testkit.TempHome(t)
	realDir := t.TempDir()
	argsFile := filepath.Join(home, "gh-args")
	testkit.FakeBinary(t, realDir, "gh", `echo "$@" > `+argsFile)

	cmd := newShimGhCmd(t, home, realDir, "pr", "list")
	// no stdin — if shim tries to prompt it will error and exit non-zero

	if err := cmd.Run(); err != nil {
		t.Fatalf("shim exited with error: %v", err)
	}

	got, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatalf("fake gh was never called — args file not found: %v", err)
	}
	if want := "pr list"; strings.TrimSpace(string(got)) != want {
		t.Errorf("fake gh args = %q, want %q", strings.TrimSpace(string(got)), want)
	}
}

func TestGhShim_PrCreate_Prompts(t *testing.T) {
	home := testkit.TempHome(t)
	realDir := t.TempDir()
	argsFile := filepath.Join(home, "gh-args")
	testkit.FakeBinary(t, realDir, "gh", `echo "$@" > `+argsFile)

	cmd := newShimGhCmd(t, home, realDir, "pr", "create")
	cmd.Stdin = strings.NewReader("y")

	if err := cmd.Run(); err != nil {
		t.Fatalf("shim exited with error: %v", err)
	}

	got, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatalf("fake gh was never called — args file not found: %v", err)
	}
	if want := "pr create"; strings.TrimSpace(string(got)) != want {
		t.Errorf("fake gh args = %q, want %q", strings.TrimSpace(string(got)), want)
	}
}
