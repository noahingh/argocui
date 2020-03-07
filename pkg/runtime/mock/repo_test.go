package mock

import (
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestNewAnimal(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *Animal
	}{
		{
			name: "create a new animal",
			args: args{
				name: "foo",
			},
			want: &Animal{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Mock",
					APIVersion: "argocui.github.com/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: namespace,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAnimal(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAnimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		r       *Repo
		args    args
		want    runtime.Object
		wantErr bool
	}{
		{
			name: "get cheetah",
			r:    &Repo{},
			args: args{
				key: "default/cheetah",
			},
			want:    NewAnimal("cheetah"),
			wantErr: false,
		},
		{
			name: "not exist",
			r:    &Repo{},
			args: args{
				key: "default/zibra",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{}
			got, err := r.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repo.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.want == nil {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repo.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_Search(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		r    *Repo
		args args
		want []runtime.Object
	}{
		{
			name: "search cheet",
			r:    &Repo{},
			args: args{
				pattern: "cheet",
			},
			want: []runtime.Object{
				NewAnimal("cheetah"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{}
			if got := r.Search(tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repo.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		r       *Repo
		args    args
		wantErr bool
	}{
		{
			name: "delete cheetah",
			r:    &Repo{},
			args: args{
				key: "default/cheetah",
			},
			wantErr: false,
		},
		{
			name: "delete non-exist",
			r:    &Repo{},
			args: args{
				key: "default/zibra",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{}
			if err := r.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Repo.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
