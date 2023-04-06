package diagnose

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"sort"
	"strings"

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

	FullObject    bool
	IncludeEvents bool
	ShowPrompt    bool
}

func (o *Options) AddFlags(flags *pflag.FlagSet) {
	o.common.AddFlags(flags)

	flags.BoolVar(&o.FullObject, "full-object", false, "Include the full object in prompt. Note that the GPT 3.5 modal can only handle 4097 tokens, including the full object might exceed the limit.")
	flags.BoolVar(&o.IncludeEvents, "include-events", true, "Include related events in prompt.")
	flags.BoolVar(&o.ShowPrompt, "show-prompt", false, "Print out the complete prompt before sending it out.")
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

	relatedEvts, err := o.getRelatedEvents(obj)
	if err != nil {
		return err
	}
	objYAML := o.marshalObject(obj)

	var buf bytes.Buffer
	_ = promptDiagnose.Execute(&buf, templateData{
		// TODO: include logs
		Data: objYAML + relatedEvts,
		Lang: o.common.Lang,
	})

	if o.ShowPrompt {
		cmd.Println(buf.String())
	}

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

func (o *Options) getRelatedEvents(obj runtime.Object) (string, error) {
	if !o.IncludeEvents {
		return "", nil
	}

	client, err := o.common.NewKubeClientSet()
	if err != nil {
		return "", err
	}
	ns, _ := meta.NewAccessor().Namespace(obj)
	eventList, err := client.CoreV1().Events(ns).Search(scheme.Scheme, obj)
	if err != nil {
		return "", err
	}
	if len(eventList.Items) == 0 {
		return "", nil
	}
	sort.Slice(eventList.Items, func(i, j int) bool {
		return eventList.Items[i].LastTimestamp.Time.Before(eventList.Items[j].LastTimestamp.Time)
	})

	var buf bytes.Buffer
	buf.WriteString("\n\n---\n\nRelated events in CSV format:\n")
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"Type", "Reason", "Age", "From", "Message"})
	for _, e := range eventList.Items {
		var interval string
		firstTimestampSince := translateMicroTimestampSince(e.EventTime)
		if e.EventTime.IsZero() {
			firstTimestampSince = translateTimestampSince(e.FirstTimestamp)
		}
		if e.Series != nil {
			interval = fmt.Sprintf("%s (x%d over %s)", translateMicroTimestampSince(e.Series.LastObservedTime), e.Series.Count, firstTimestampSince)
		} else if e.Count > 1 {
			interval = fmt.Sprintf("%s (x%d over %s)", translateTimestampSince(e.LastTimestamp), e.Count, firstTimestampSince)
		} else {
			interval = firstTimestampSince
		}
		source := e.Source.Component
		if source == "" {
			source = e.ReportingController
		}
		_ = w.Write([]string{
			e.Type,
			e.Reason,
			interval,
			source,
			strings.TrimSpace(e.Message),
		})
	}
	w.Flush()

	return buf.String(), nil
}
