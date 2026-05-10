package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/rockwellwindsor/tollgate/internal/state"
)

func newPauseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pause",
		Short: "Disable prompts for this shell session",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := state.NewManager(state.DefaultDir())
			if err := mgr.SetSessionPaused(os.Getpid()); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "tollgate: session paused")
			return nil
		},
	}
}

func newResumeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resume",
		Short: "Re-enable prompts for this shell session",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := state.NewManager(state.DefaultDir())
			if err := mgr.ClearSessionPaused(os.Getpid()); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "tollgate: session resumed")
			return nil
		},
	}
}

func newOnCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "on",
		Short: "Enable tollgate globally",
		RunE:  func(cmd *cobra.Command, args []string) error { return nil },
	}
}

func newOffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "off",
		Short: "Disable tollgate globally",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := state.NewManager(state.DefaultDir())
			if err := mgr.SetGlobalOff(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "tollgate: disabled")
			return nil
		},
	}
}
