package serializer

import (
	"reflect"
	"testing"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestConvertToNamespace(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    runtime.Object
		wantErr bool
	}{
		{
			name: "default namespace",
			args: args{
				data: []byte(`apiVersion: v1
kind: Namespace
metadata:
  name: default
spec:
  finalizers:
  - kubernetes
status:
  phase: Active
`),
			},
			want: &corev1.Namespace{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Namespace",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
				Spec: corev1.NamespaceSpec{
					Finalizers: []corev1.FinalizerName{"kubernetes"},
				},
				Status: corev1.NamespaceStatus{
					Phase: corev1.NamespaceActive,
				},
			},
			wantErr: false,
		},
		{
			name: "wrong GroupVersionKind",
			args: args{
				data: []byte(`apiVersion: v1
kind: Foo
metadata:
  name: foo
`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToNamespace(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToWorkflow(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    runtime.Object
		wantErr bool
	}{
		{
			name: "hello world",
			args: args{
				data: []byte(`apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: hello-world-
spec:
  entrypoint: whalesay
  templates:
  - name: whalesay
    container:
      image: docker/whalesay:latest
      command: [cowsay]
      args: ["hello world"]`),
			},
			want: &wf.Workflow{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "argoproj.io/v1alpha1",
					Kind: "Workflow",
				},
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "hello-world-",
				},
				Spec: wf.WorkflowSpec{
					Entrypoint: "whalesay",
					Templates: []wf.Template{
						wf.Template{
							Name: "whalesay",
							Container: &corev1.Container{
								Image: "docker/whalesay:latest",
								Command: []string{"cowsay"},
								Args: []string{"hello world"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToWorkflow(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToWorkflow() = %v, want %v", got, tt.want)
			}
		})
	}
}
