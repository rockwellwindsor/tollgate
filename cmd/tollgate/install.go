package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/rockwellwindsor/tollgate/internal/state"
)

func newInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Write shim binaries and config",
		RunE: func(cmd *cobra.Command, args []string) error {
			exe, err := os.Executable()
			if err != nil {
				return err
			}
			return install(cmd, filepath.Dir(exe))
		},
	}
}

func install(cmd *cobra.Command, srcDir string) error {
	home := state.DefaultDir()
	binDir := filepath.Join(home, "bin")

	if err := os.MkdirAll(binDir, 0755); err != nil {
		return err
	}

	shims := map[string]string{
		"shim-git": "git",
		"shim-gh":  "gh",
	}
	for src, dst := range shims {
		if err := copyExe(filepath.Join(srcDir, src), filepath.Join(binDir, dst)); err != nil {
			return err
		}
	}

	if err := copyExe(filepath.Join(srcDir, "tollgate"), filepath.Join(binDir, "tollgate")); err != nil {
		return err
	}

	configPath := filepath.Join(home, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.WriteFile(configPath, []byte("{}\n"), 0644); err != nil {
			return err
		}
	}

	fmt.Fprintf(cmd.OutOrStdout(), "installed to %s\n", binDir)
	fmt.Fprintf(cmd.OutOrStdout(), "add to your shell profile:\n  export PATH=\"%s:$PATH\"\n", binDir)
	return nil
}

func copyExe(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
