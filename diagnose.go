package main

import (
	"github.com/spf13/cobra"
)

func newDiagnoseCommand() *cobra.Command {
	return &cobra.Command{
		Use: "diagnose",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}
