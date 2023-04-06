package diagnose

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"

	"github.com/knight42/kopilot/cmd"
)

type Options struct {
	common cmd.Options

	FullObject bool
}

func (o *Options) AddFlags(flags *pflag.FlagSet) {
	o.common.AddFlags(flags)

	flags.BoolVar(&o.FullObject, "full-object", false, "Include the full object in prompt. Note that the GPT 3.5 modal can only handle 4097 tokens, including the full object might exceed the limit.")
}

func (o *Options) Run(cmd *cobra.Command, args []string) error {
	ns, _, _ := o.common.ToRawKubeConfigLoader().Namespace()
	obj, err := o.common.NewBuilder().
		NamespaceParam(ns).
		DefaultNamespace().
		WithScheme(scheme.Scheme, scheme.Scheme.PrioritizedVersionsAllGroups()...).
		ResourceNames(args[0], args[1]).
		Do().
		Object()
	if err != nil {
		return fmt.Errorf("get object: %w", err)
	}

	data := o.marshalObject(obj)
	var buf bytes.Buffer
	_ = promptDiagnose.Execute(&buf, templateData{
		// TODO: include logs or events
		Data: data,
		Lang: o.common.Lang,
	})

	return o.common.NewChatGPTClient(" Diagnosing...").
		CreateCompletion(cmd.Context(), buf.String(), cmd.OutOrStdout())
}

func (o *Options) marshalObject(obj runtime.Object) string {
	if a, err := meta.Accessor(obj); err == nil {
		a.SetManagedFields(nil)
		a.SetAnnotations(nil)
		a.SetLabels(nil)
	}

	switch actual := obj.(type) {
	case *corev1.Pod:
		if !o.FullObject {
			// Usually we only need to know the pod's status.
			actual.Spec = corev1.PodSpec{}
		}
	case *corev1.Node:
		actual.Status.Images = nil
	}

	data, _ := yaml.Marshal(obj)
	return string(data)
}
