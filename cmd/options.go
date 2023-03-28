package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/knight42/kopilot/internal/client"
)

const (
	EnvKopilotToken = "KOPILOT_TOKEN"
	EnvKopilotLang  = "KOPILOT_LANG"
)

func getEnvOrDefault(key, defVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defVal
}

type Options struct {
	// set by env vars
	Token string
	Lang  string

	// set by flags
	kubeConfigFlags *genericclioptions.ConfigFlags
}

func NewOptions() Options {
	return Options{
		kubeConfigFlags: genericclioptions.NewConfigFlags(true),
	}
}

func (o *Options) Complete() error {
	o.Token = os.Getenv(EnvKopilotToken)
	if len(o.Token) == 0 {
		return fmt.Errorf("please specify the api token with ENV %s", EnvKopilotToken)
	}
	o.Lang = getEnvOrDefault(EnvKopilotLang, "English")
	return nil
}

func (o *Options) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(o.kubeConfigFlags.KubeConfig, "kubeconfig", *o.kubeConfigFlags.KubeConfig, "Path to the kubeconfig file to use for CLI requests.")
	flags.StringVarP(o.kubeConfigFlags.Namespace, "namespace", "n", *o.kubeConfigFlags.Namespace, "If present, the namespace scope for this CLI request")
}

func (o *Options) NewBuilder() *resource.Builder {
	return resource.NewBuilder(o.kubeConfigFlags)
}

func (o *Options) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return o.kubeConfigFlags.ToRawKubeConfigLoader()
}

func (o *Options) NewChatGPTClient(spinnerSuffix string) client.Client {
	return client.NewChatGPTClient(o.Token, spinnerSuffix)
}

func (o *Options) NewKubeClientSet() (kubernetes.Interface, error) {
	cfg, err := o.kubeConfigFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("create rest config: %w", err)
	}
	return kubernetes.NewForConfig(cfg)
}
