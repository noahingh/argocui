package mock

import (
	"context"
	"reflect"
	"testing"

	"github.com/hanjunlee/argocui/pkg/resource"
)

func TestRepo_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		r    *Repo
		args args
		want interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{}
			if got := r.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
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
		want []interface{}
	}{
		// TODO: Add test cases.
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
			name: "nothing",
			r: &Repo{},
			args: args{
				key: "zibra",
			},
			wantErr: true,
		},
		{
			name: "rm end element",
			r: &Repo{},
			args: args{
				key: "cheetah",
			},
			wantErr: false,
		},
		{
			name: "rm camel",
			r: &Repo{},
			args: args{
				key: "camel",
			},
			wantErr: false,
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

func TestRepo_Logs(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		r       *Repo
		args    args
		wantCh  <-chan resource.Log
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{}
			gotCh, err := r.Logs(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repo.Logs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCh, tt.wantCh) {
				t.Errorf("Repo.Logs() = %v, want %v", gotCh, tt.wantCh)
			}
		})
	}
}
