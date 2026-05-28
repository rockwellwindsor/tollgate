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
			paused, err := mgr.IsSessionPaused(os.Getpid())
			if err != nil {
				return err
			}
			if paused {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "tollgate: already paused")
				return nil
			}
			if err := mgr.SetSessionPaused(os.Getpid()); err != nil {
				return err
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "tollgate: session paused")
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
			paused, err := mgr.IsSessionPaused(os.Getpid())
			if err != nil {
				return err
			}
			if !paused {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "tollgate: not paused")
				return nil
			}
			if err := mgr.ClearSessionPaused(os.Getpid()); err != nil {
				return err
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "tollgate: session resumed")
			return nil
		},
	}
}

func newOnCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "on",
		Short: "Enable tollgate globally",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := state.NewManager(state.DefaultDir())
			on, err := mgr.IsGlobalOn()
			if err != nil {
				return err
			}
			if on {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "tollgate: already enabled")
				return nil
			}
			if err := mgr.SetGlobalOn(); err != nil {
				return err
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "tollgate: enabled")
			return nil
		},
	}
}

func newOffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "off",
		Short: "Disable tollgate globally",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := state.NewManager(state.DefaultDir())
			on, err := mgr.IsGlobalOn()
			if err != nil {
				return err
			}
			if !on {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "tollgate: already disabled")
				return nil
			}
			if err := mgr.SetGlobalOff(); err != nil {
				return err
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "tollgate: disabled")
			return nil
		},
	}
}
