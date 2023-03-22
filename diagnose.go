package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"sigs.k8s.io/yaml"

	"github.com/knight42/kopilot/client"
)

const (
	envKopilotToken = "KOPILOT_TOKEN"
	envKopilotType  = "KOPILOT_TOKEN_TYPE"
	envKopilotLang  = "KOPILOT_LANG"

	typeChatGPT = "ChatGPT"
	langEN      = "English"
	langCN      = "Chinese"
)

type option struct {
	token string
	lang  string
	typ   string
}

type templateData struct {
	Data string
	Lang string
}

func newDiagnoseCommand(opt option) *cobra.Command {
	cf := genericclioptions.NewConfigFlags(true)
	cmd := &cobra.Command{
		Use:   "diagnose TYPE NAME",
		Short: "Diagnose a resource",
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
			_ = promptDiagnose.Execute(&buf, templateData{
				Data: data,
				Lang: opt.lang,
			})

			cmd.Println("Diagnosing...")
			resp, err := createCompletion(cmd.Context(), opt, buf.String())
			if err != nil {
				return err
			}
			cmd.Println(resp)
			return nil
		},
	}
	flags := cmd.Flags()
	flags.StringVar(cf.KubeConfig, "kubeconfig", *cf.KubeConfig, "Path to the kubeconfig file to use for CLI requests.")
	flags.StringVarP(cf.Namespace, "namespace", "n", *cf.Namespace, "If present, the namespace scope for this CLI request")

	return cmd
}

func createCompletion(ctx context.Context, opt option, prompt string) (string, error) {
	var cli client.Client
	switch opt.typ {
	case typeChatGPT:
		cli = client.NewChatGPTClient(opt.token)
	default:
		return "", fmt.Errorf("invalid type %s", opt.typ)
	}
	resp, err := cli.CreateCompletion(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("create completion: %w", err)
	}
	return resp, nil
}

func getResourceInYAML(cf *genericclioptions.ConfigFlags, ns, res, name string) (string, error) {
	obj, err := resource.NewBuilder(cf).
		NamespaceParam(ns).
		DefaultNamespace().
		Unstructured().
		SingleResourceType().
		ResourceNames(res, name).
		Do().
		Object()
	if err != nil {
		return "", fmt.Errorf("get object: %w", err)
	}
	metaObj, err := meta.Accessor(obj)
	if err == nil {
		metaObj.SetManagedFields(nil)
	}

	data, err := yaml.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("marshal object: %w", err)
	}
	return string(data), nil
}
