package resolver

import (
	"fmt"
	"os"
	"path/filepath"
)

func ResolveRealBinary(name, skipDir string) (string, error) {
	skip := cleanPath(skipDir)

	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		if cleanPath(dir) == skip {
			continue
		}
		candidate := filepath.Join(dir, name)
		if isExecutable(candidate) {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("no real %q found on PATH", name)
}

func cleanPath(p string) string {
	return filepath.Clean(p)
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0111 != 0
}
