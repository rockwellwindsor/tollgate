package testkit

import (
	"os"
	"os/exec"
	"strings"
)

type ShimResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// RunShim executes the shim binary at shimPath with the given args and PATH.
// extraEnv entries are appended after the inherited environment.
func RunShim(shimPath string, args []string, extraEnv []string) ShimResult {
	cmd := exec.Command(shimPath, args...)
	cmd.Env = append(os.Environ(), extraEnv...)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := ShimResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		result.ExitCode = exitErr.ExitCode()
	}
	return result
}
