package ui

import (
	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	scheme.AddKnownTypes(schema.GroupVersion{Group: "", Version: "v1"}, &corev1.Namespace{})
	scheme.AddKnownTypes(schema.GroupVersion{Group: "argoproj.io", Version: "v1alpha1"}, &wf.Workflow{})
}

func objectKind(o runtime.Object) (schema.GroupVersionKind, bool, error) {
	gvks, isUnversionedType, err := scheme.ObjectKinds(o)

	if len(gvks) == 0 {
		return schema.GroupVersionKind{}, isUnversionedType, err
	}

	return gvks[0], isUnversionedType, err
}
