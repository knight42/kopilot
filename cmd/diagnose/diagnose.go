package diagnose

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"

	"github.com/knight42/kopilot/cmd"
)

func New(commonOpts cmd.Options) *cobra.Command {
	c := &cobra.Command{
		Use:          "diagnose TYPE NAME",
		Short:        "Diagnose a resource",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ns, _, _ := commonOpts.ToRawKubeConfigLoader().Namespace()
			obj, err := commonOpts.NewBuilder().
				NamespaceParam(ns).
				DefaultNamespace().
				WithScheme(scheme.Scheme, scheme.Scheme.PrioritizedVersionsAllGroups()...).
				ResourceNames(args[0], args[1]).
				Do().
				Object()
			if err != nil {
				return fmt.Errorf("get object: %w", err)
			}

			pruneObject(obj)

			// TODO: include logs or events
			data, err := yaml.Marshal(obj)
			if err != nil {
				return fmt.Errorf("marshal object: %w", err)
			}

			var buf bytes.Buffer
			_ = promptDiagnose.Execute(&buf, templateData{
				Data: string(data),
				Lang: commonOpts.Lang,
			})

			// TODO: limit the length of the prompt?
			return commonOpts.NewChatGPTClient(" Diagnosing...").
				CreateCompletion(cmd.Context(), buf.String(), cmd.OutOrStdout())
		},
	}
	flags := c.Flags()
	commonOpts.AddFlags(flags)
	return c
}
