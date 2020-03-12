package ui

import (
	"github.com/hanjunlee/argocui/pkg/runtime/mock"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	scheme.AddKnownTypes(schema.GroupVersion{Group: "argocui.github.com", Version: "v1"}, &mock.Animal{})
	scheme.AddKnownTypes(schema.GroupVersion{Group: "", Version: "v1"}, &corev1.Namespace{})
	// TODO: add the scheme workflow.
}

func objectKind(o runtime.Object) (schema.GroupVersionKind, bool, error) {
	gvks, isUnversionedType, err := scheme.ObjectKinds(o)
	
	if len(gvks) == 0 {
		return schema.GroupVersionKind{}, isUnversionedType, err
	}

	return gvks[0], isUnversionedType, err
}
