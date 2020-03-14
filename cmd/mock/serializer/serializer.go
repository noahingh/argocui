package serializer

import (
	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

var (
	scheme     = runtime.NewScheme()
	serializer = json.NewSerializerWithOptions(json.SimpleMetaFactory{}, scheme, scheme, json.SerializerOptions{
		Yaml: true,
	})
)

func init() {
	scheme.AddKnownTypes(schema.GroupVersion{Group: "", Version: "v1"}, &corev1.Namespace{})
	scheme.AddKnownTypes(schema.GroupVersion{Group: "argoproj.io", Version: "v1alpha1"}, &wf.Workflow{})
}


// ConvertToNamespace converts a YAML into a object.
func ConvertToNamespace(data []byte) (runtime.Object, error) {
	n, _, err := serializer.Decode(data, &schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"}, &corev1.Namespace{})
	return n, err
}

// ConvertToWorkflow converts a YAML into a object.
func ConvertToWorkflow(data []byte) (runtime.Object, error) {
	w, _, err := serializer.Decode(data, &schema.GroupVersionKind{Group: "argoproj.io", Version: "v1alpha1", Kind: "Workflow"}, &wf.Workflow{})
	return w, err
}
