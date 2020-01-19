package argo

import (
	"context"
	"time"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
)

// UseCase is the interface of use case.
type UseCase interface {
	Get(key string) *wf.Workflow
	Search(pattern string) []*wf.Workflow
	Delete(key string) error
	// Logs get the channel to recieve Logs from a Argo workflow.
	Logs(ctx context.Context, key string) (ch <-chan Log, err error)
}

// Repository is the interface of repositories.
type Repository interface {
	Get(key string) *wf.Workflow
	Search(pattern string) []*wf.Workflow
	Delete(key string) error
	// Logs get the channel to recieve Logs from a Argo workflow.
	Logs(ctx context.Context, key string) (ch <-chan Log, err error)
}

// Log  is log from a Argo workflow.
type Log struct {
	Pod         string
	Message     string
	Time        time.Time
}
