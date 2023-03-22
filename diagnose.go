package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"sigs.k8s.io/yaml"
)

type templateData struct {
	Data string
}

func newDiagnoseCommand() *cobra.Command {
	cf := genericclioptions.NewConfigFlags(true)
	cmd := &cobra.Command{
		Use:          "diagnose TYPE NAME",
		Short:        "Diagnose a resource",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("must specify type and name")
			}
			ns, _, _ := cf.ToRawKubeConfigLoader().Namespace()
			obj, err := resource.NewBuilder(cf).
				NamespaceParam(ns).
				DefaultNamespace().
				Unstructured().
				SingleResourceType().
				ResourceNames(args[0], args[1]).
				Do().
				Object()
			if err != nil {
				return fmt.Errorf("get object: %w", err)
			}
			metaObj, err := meta.Accessor(obj)
			if err == nil {
				metaObj.SetManagedFields(nil)
			}

			data, err := yaml.Marshal(obj)
			if err != nil {
				return fmt.Errorf("marshal object: %w", err)
			}

			return promptDiagnose.Execute(cmd.OutOrStdout(), templateData{
				Data: string(data),
			})
		},
	}
	flags := cmd.Flags()
	flags.StringVar(cf.KubeConfig, "kubeconfig", *cf.KubeConfig, "Path to the kubeconfig file to use for CLI requests.")
	flags.StringVarP(cf.Namespace, "namespace", "n", *cf.Namespace, "If present, the namespace scope for this CLI request")

	return cmd
}
