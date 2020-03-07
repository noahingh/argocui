package runtime

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// UseCase is the interface of use case.
type UseCase interface {
	Get(key string) (runtime.Object, error)
	Search(pattern string) []runtime.Object
	Delete(key string) error
}

// Repo is the interface of repository.
type Repo interface {
	Get(key string) (runtime.Object, error)
	Search(pattern string) []runtime.Object
	Delete(key string) error
}
