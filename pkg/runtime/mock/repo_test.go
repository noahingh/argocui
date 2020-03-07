package mock

import (
	"reflect"
	"testing"
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
		// TODO: Add test cases.
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
		want    *Animal
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{}
			got, err := r.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repo.Get() error = %v, wantErr %v", err, tt.wantErr)
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
		want []*Animal
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
		// TODO: Add test cases.
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
