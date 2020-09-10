package namespace

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestPresentor_convertToRows(t *testing.T) {
	type args struct {
		objs []runtime.Object
	}
	tests := []struct {
		name string
		p    *Presentor
		args args
		want [][]string
	}{
		{
			name: "namespace objects",
			p: &Presentor{},
			args: args{
				objs: []runtime.Object{
					&corev1.Namespace{
						ObjectMeta: metav1.ObjectMeta{
							Name: "default",
						},
					},
				},
			},
			want: [][]string{
				{"default"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Presentor{}
			if got := p.convertToRows(tt.args.objs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Presentor.convertToRows() = %v, want %v", got, tt.want)
			}
		})
	}
}
