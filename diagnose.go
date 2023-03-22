package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"sigs.k8s.io/yaml"
)

const (
	envKopilotToken = "KOPILOT_TOKEN"
	envKopilotType  = "KOPILOT_TOKEN_TYPE"
	envKopilotLang  = "KOPILOT_LANG"

	typeChatGPT = "ChatGPT"
	langEN      = "English"
	langCN      = "Chinese"
)

func preCheck(cmd *cobra.Command, args []string) error {
	token := os.Getenv(envKopilotToken)
	typ := os.Getenv(envKopilotType)
	if typ == "" {
		typ = typeChatGPT
	}
	lang := os.Getenv(envKopilotLang)
	if lang == "" {
		lang = langEN
	}
	cmd.Printf("You're using %s client with %s, customize them with %s and %s.\n", typ, lang, envKopilotType, envKopilotLang)
	if token == "" {
		return fmt.Errorf("please specify the token for %s, please specify the token with ENV %s", typ, envKopilotToken)
	}
	return nil
}

type templateData struct {
	Data string
}

func newDiagnoseCommand() *cobra.Command {
	cf := genericclioptions.NewConfigFlags(true)
	cmd := &cobra.Command{
		Use:          "diagnose TYPE NAME",
		Short:        "Diagnose a resource",
		PreRunE:      preCheck,
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
