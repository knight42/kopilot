package diagnose

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

func pruneObject(obj runtime.Object) {
	if a, err := meta.Accessor(obj); err == nil {
		a.SetManagedFields(nil)
		a.SetAnnotations(nil)
		a.SetLabels(nil)
	}

	switch actual := obj.(type) {
	case *corev1.Node:
		pruneNode(actual)
	}
}

func pruneNode(obj *corev1.Node) {
	obj.Status.Images = nil
}
