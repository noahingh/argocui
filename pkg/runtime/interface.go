package runtime

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
)

// Log  is the struct of log.
type Log struct {
	DisplayName string
	Pod         string
	Message     string
	Time        time.Time
}

// UseCase is the interface of use case.
type UseCase interface {
	Get(key string) (runtime.Object, error)
	Search(namespace, pattern string) []runtime.Object
	Delete(key string) error
	Logs(ctx context.Context, key string) (<-chan Log, error)
}

// Repo is the interface of repository.
type Repo interface {
	Get(key string) (runtime.Object, error)
	Search(namespace, pattern string) []runtime.Object
	Delete(key string) error
	Logs(ctx context.Context, key string) (<-chan Log, error)
}
