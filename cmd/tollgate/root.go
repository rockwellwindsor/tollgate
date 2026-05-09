package main

import "github.com/spf13/cobra"

const rootLong = `tollgate is a checkpoint for sensitive operations.

Place shim binaries on your PATH ahead of the real git and gh.
Watched commands (git push, gh pr create, etc.) pause for confirmation.`

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "tollgate",
		Short:   "A checkpoint for sensitive operations",
		Long:    rootLong,
		Version: "0.1.0",
	}

	root.AddCommand(
		&cobra.Command{Use: "status", Short: "Show current state and active rules"},
		&cobra.Command{Use: "pause", Short: "Disable prompts for this shell session"},
		&cobra.Command{Use: "resume", Short: "Re-enable prompts for this shell session"},
		&cobra.Command{Use: "on", Short: "Enable tollgate globally"},
		&cobra.Command{Use: "off", Short: "Disable tollgate globally"},
		&cobra.Command{Use: "clear-logs", Short: "Delete the audit log"},
		&cobra.Command{Use: "install", Short: "Write shim binaries and config"},
		&cobra.Command{Use: "config", Short: "Show or edit configuration"},
	)

	return root
}
