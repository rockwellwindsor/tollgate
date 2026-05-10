package main

import (
	"fmt"

	"github.com/spf13/cobra"

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
	mgr := state.NewManager(state.DefaultDir())

	on, err := mgr.IsGlobalOn()
	if err != nil {
		return err
	}
	if on {
		fmt.Fprintln(cmd.OutOrStdout(), "tollgate: ENABLED")
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), "tollgate: DISABLED")
	}

	return nil
}
