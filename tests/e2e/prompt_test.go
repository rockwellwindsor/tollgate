package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rockwellwindsor/tollgate/internal/testkit"
)

func TestPrompt_YesInvokesFakeGit(t *testing.T) {
	home := testkit.TempHome(t)
	t.Setenv("TOLLGATE_HOME", home)

	realDir := t.TempDir()
	argsFile := filepath.Join(home, "git-args")
	testkit.FakeBinary(t, realDir, "git", `echo "$@" > `+argsFile)

	shimDir := filepath.Dir(shimGitPath)
	cmd := exec.Command(shimGitPath, "push", "origin", "main")
	cmd.Env = append(os.Environ(),
		"PATH="+shimDir+string(os.PathListSeparator)+realDir,
		"HOME="+home,
		"TOLLGATE_HOME="+home,
	)
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
