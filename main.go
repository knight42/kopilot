package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/knight42/kopilot/cmd"
	"github.com/knight42/kopilot/cmd/audit"
	"github.com/knight42/kopilot/cmd/diagnose"
)

func getCommandName() string {
	if strings.HasPrefix(filepath.Base(os.Args[0]), "kubectl-") {
		// cobra will split on " " and take the first element
		return "kubectl\u2002kopilot"
	}
	return "kopilot"
}

func main() {
	opt := cmd.NewOptions()
	err := opt.Complete()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	rootCmd := cobra.Command{
		Use: getCommandName(),
		Long: fmt.Sprintf(`You need to set TWO ENVs to run Kopilot.
Set %s to specify your token.
Set %s to specify the language. Valid options like Chinese, French, Spain, etc.
`, cmd.EnvKopilotToken, cmd.EnvKopilotLang),
	}
	rootCmd.AddCommand(
		audit.New(opt),
		diagnose.New(opt),
	)
	_ = rootCmd.Execute()
}
