package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/rockwellwindsor/tollgate/internal/state"
)

func newClearLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear-logs",
		Short: "Delete the audit log",
		RunE:  runClearLogsCmd,
	}
	cmd.Flags().Bool("dry-run", false, "Show what would be deleted without deleting")
	cmd.Flags().Bool("yes", false, "Skip confirmation prompt")
	return cmd
}

func runClearLogsCmd(cmd *cobra.Command, _ []string) error {
	logPath := filepath.Join(state.DefaultDir(), "audit.log")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	yes, _ := cmd.Flags().GetBool("yes")

	if dryRun {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "would delete: %s\n", logPath)
		return nil
	}

	if yes {
		return deleteLog(logPath)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "delete %s? [y/N] ", logPath)
	buf := make([]byte, 1)
	if _, err := cmd.InOrStdin().Read(buf); err != nil {
		return err
	}
	if buf[0] == 'y' || buf[0] == 'Y' {
		return deleteLog(logPath)
	}
	return nil
}

func deleteLog(path string) error {
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
