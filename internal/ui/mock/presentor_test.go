package mock

import (
	"reflect"
	"testing"

	mockrepo "github.com/hanjunlee/argocui/pkg/runtime/mock"

	"k8s.io/apimachinery/pkg/runtime"
)

func NewObjects(ss... string) []runtime.Object {
	objs := make([]runtime.Object, 0)
	for _, s := range ss {
		a := mockrepo.NewAnimal(s)
		objs = append(objs, a)
	}
	return objs
}

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
			name: "convert",
			p: &Presentor{},
			args: args{
				objs: NewObjects("alligator", "ant"),
			},
			want: [][]string{
				[]string{"default", "alligator"},
				[]string{"default", "ant"},
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
