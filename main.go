package main

import (
	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use: "kopilot",
	}
	cmd.AddCommand(newDiagnoseCommand())
	_ = cmd.Execute()
}
