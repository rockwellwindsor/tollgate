package testkit

import "testing"

// TempHome redirects $HOME to a fresh temp directory for the duration of the test.
func TempHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	return dir
}
