package main

import "github.com/spf13/cobra"

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "tollgate",
		Short:   "A checkpoint for sensitive operations",
		Version: "0.1.0",
	}
}
