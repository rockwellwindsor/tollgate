package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func newShimCmdFor(t *testing.T, shimPath, home, realDir string, args ...string) *exec.Cmd {
	t.Helper()
	shimDir := filepath.Dir(shimPath)
	cmd := exec.Command(shimPath, args...)
	cmd.Env = append(os.Environ(),
		"PATH="+shimDir+string(os.PathListSeparator)+realDir,
		"HOME="+home,
		"TOLLGATE_HOME="+home,
	)
	return cmd
}
