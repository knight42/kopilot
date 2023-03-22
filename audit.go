package main

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func newAuditCommand(opt option) *cobra.Command {
	cf := genericclioptions.NewConfigFlags(true)
	cmd := &cobra.Command{
		Use:   "audit TYPE NAME",
		Short: "Audit a resource",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opt.token == "" {
				return fmt.Errorf("please specify the token for %s with ENV %s", opt.typ, envKopilotToken)
			}
			return nil
		},
		SilenceUsage: true,
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ns, _, _ := cf.ToRawKubeConfigLoader().Namespace()
			data, err := getResourceInYAML(cf, ns, args[0], args[1])
			if err != nil {
				return fmt.Errorf("marshal object: %w", err)
			}

			var buf bytes.Buffer
			_ = promptAudit.Execute(&buf, templateData{
				Data: data,
				Lang: opt.lang,
			})

			cmd.Println("Auditing...")
			err = createCompletion(cmd.Context(), opt, buf.String(), cmd.OutOrStderr())
			if err != nil {
				return err
			}
			return nil
		},
	}
	flags := cmd.Flags()
	flags.StringVar(cf.KubeConfig, "kubeconfig", *cf.KubeConfig, "Path to the kubeconfig file to use for CLI requests.")
	flags.StringVarP(cf.Namespace, "namespace", "n", *cf.Namespace, "If present, the namespace scope for this CLI request")

	return cmd
}
