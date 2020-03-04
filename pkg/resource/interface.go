package resource

import (
	"context"
	"time"
)

// UseCase is the interface of use case.
type UseCase interface {
	Get(key string) interface{}
	Search(pattern string) []interface{}
	Delete(key string) error
	Logs(ctx context.Context, key string) (ch <-chan Log, err error)
}

// Repo is the interface of repository.
type Repo interface {
	Get(key string) interface{}
	Search(pattern string) []interface{}
	Delete(key string) error
	Logs(ctx context.Context, key string) (ch <-chan Log, err error)
}

// Log  is log from a Argo workflow.
type Log struct {
	DisplayName string
	Pod         string
	Message     string
	Time        time.Time
}
