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
		Version: "0.1.1",
	}

	root.AddCommand(
		newStatusCmd(),
		newPauseCmd(),
		newResumeCmd(),
		newOnCmd(),
		newOffCmd(),
		newClearLogsCmd(),
		newInstallCmd(),
		&cobra.Command{Use: "config", Short: "Show or edit configuration"},
	)

	return root
}
