package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/rockwellwindsor/tollgate/internal/audit"
	"github.com/rockwellwindsor/tollgate/internal/state"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current state and active rules",
		RunE:  runStatusCmd,
	}
}

func runStatusCmd(cmd *cobra.Command, _ []string) error {
	dir := state.DefaultDir()
	mgr := state.NewManager(dir)

	on, err := mgr.IsGlobalOn()
	if err != nil {
		return err
	}
	if on {
		fmt.Fprintln(cmd.OutOrStdout(), "tollgate: ENABLED")
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), "tollgate: DISABLED")
	}

	paused, err := mgr.IsSessionPaused(os.Getpid())
	if err != nil {
		return err
	}
	if paused {
		fmt.Fprintln(cmd.OutOrStdout(), "session: paused")
	}

	patterns, err := mgr.AllowedPatterns(os.Getpid())
	if err != nil {
		return err
	}
	for _, p := range patterns {
		fmt.Fprintf(cmd.OutOrStdout(), "session allow: %s\n", p)
	}

	logPath := filepath.Join(dir, "audit.log")
	entries, err := audit.Read(logPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "audit: %d entries\n", len(entries))

	return nil
}
