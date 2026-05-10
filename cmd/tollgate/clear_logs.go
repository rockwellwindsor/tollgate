package main

import (
	"fmt"
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

	if dryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "would delete: %s\n", logPath)
		return nil
	}

	return nil
}
