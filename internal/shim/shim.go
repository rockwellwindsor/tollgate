package shim

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/rockwellwindsor/tollgate/internal/resolver"
)

func Run(binaryName string) {
	self, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "tollgate: could not determine shim path: %v\n", err)
		os.Exit(1)
	}

	real, err := resolver.ResolveRealBinary(binaryName, filepath.Dir(self))
	if err != nil {
		fmt.Fprintf(os.Stderr, "tollgate: %v\n", err)
		os.Exit(1)
	}

	if err := syscall.Exec(real, append([]string{real}, os.Args[1:]...), os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "tollgate: exec failed: %v\n", err)
		os.Exit(1)
	}
}
