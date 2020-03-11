package ui

import (
	"reflect"
	"testing"

	"github.com/hanjunlee/argocui/pkg/runtime/mock"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Test_objectKind(t *testing.T) {
	type args struct {
		o runtime.Object
	}
	tests := []struct {
		name    string
		args    args
		want    schema.GroupVersionKind
		want1   bool
		wantErr bool
	}{
		{
			name: "kind Animal",
			args: args{
				o: mock.NewAnimal("alligator"),
			},
			want: schema.GroupVersionKind{
				Group: "argocui.github.com",
				Version: "v1",
				Kind: "Animal",
			},
			want1: false,
			wantErr: false,
		},
		{
			name: "kind Namespace",
			args: args{
				o: &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "default",
					},
				},
			},
			want: schema.GroupVersionKind{
				Group: "",
				Version: "v1",
				Kind: "Namespace",
			},
			want1: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := objectKind(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("objectKind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("objectKind() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("objectKind() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
