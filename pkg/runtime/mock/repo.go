package mock

import (
	"fmt"
	"strings"
	"context"

	svc "github.com/hanjunlee/argocui/pkg/runtime"
	err "github.com/hanjunlee/argocui/pkg/runtime/error"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

const (
	namespace = "default"
)

// Animal is the name of animal.
type Animal struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

// NewAnimal create a new animal.
func NewAnimal(name string) *Animal {
	return &Animal{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Animal",
			APIVersion: "argocui.github.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

// DeepCopyObject copy the animal object.
func (a *Animal) DeepCopyObject() runtime.Object {
	return NewAnimal(a.GetName())
}

var (
	zoo = []*Animal{
		NewAnimal("alligator"), NewAnimal("ant"), NewAnimal("bear"), NewAnimal("bee"), NewAnimal("camel"), NewAnimal("cat"), NewAnimal("cheetah"),
	}
)

// Repo is the mocking repository.
type Repo struct{}

// Get return animal.
func (r *Repo) Get(key string) (runtime.Object, error) {
	for _, a := range zoo {
		ka, _ := cache.MetaNamespaceKeyFunc(a)
		if ka == key {
			return a, nil
		}
	}

	return nil, fmt.Errorf("there's not exist: %s", key)
}

// Search return animal which is matched with the pattern.
func (r *Repo) Search(namespace, pattern string) []runtime.Object {
	animals := make([]runtime.Object, 0)
	for _, a := range zoo {
		ka, _ := cache.MetaNamespaceKeyFunc(a)
		if i := strings.Index(ka, namespace+"/"+pattern); i != -1 {
			animals = append(animals, a)
		}
	}
	return animals
}

// Delete delete the animal.
func (r *Repo) Delete(key string) error {
	for idx, a := range zoo {
		ka, _ := cache.MetaNamespaceKeyFunc(a)
		if ka != key {
			continue
		}

		if idx == len(zoo)-1 {
			zoo = zoo[0:idx]
		} else {
			zoo = append(zoo[:idx], zoo[idx+1:len(zoo)]...)
		}
		return nil
	}

	return fmt.Errorf("it doesn't exist: %s", key)
}

// Logs is not implemented.
func (r *Repo) Logs(ctx context.Context, key string) (<-chan svc.Log, error) {
	return nil, err.ErrNotImplement
}
