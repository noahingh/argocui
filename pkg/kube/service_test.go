package kube

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func newMockRepo(namespaces []string) *fake.Clientset {
	objects := make([]runtime.Object, 0)
	for _, namespace := range namespaces {
		n := &corev1.Namespace{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Namespace",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		objects = append(objects, runtime.Object(n))
	}

	return fake.NewSimpleClientset(objects...)
}

func TestService_GetNamespaces(t *testing.T) {
	type fields struct {
		repo kubernetes.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "single namespace",
			fields: fields{
				repo: newMockRepo([]string{"default"}),
			},
			want: []string{"default"},
			wantErr: false,
		},
		{
			name: "multiple namespaces",
			fields: fields{
				repo: newMockRepo([]string{"default", "argo", "foo"}),
			},
			want: []string{"default", "argo", "foo"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repo: tt.fields.repo,
			}
			got, err := s.GetNamespaces()
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetNamespaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetNamespaces() = %v, want %v", got, tt.want)
			}
		})
	}
}
