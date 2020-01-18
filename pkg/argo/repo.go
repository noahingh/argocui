package argo

import (
	"context"
	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
)

// Repository is the interface of repositories.
type Repository interface {
	Get(key string) *wf.Workflow
	Search(pattern string) []*wf.Workflow
	Delete(key string) error
	Logs(ctx context.Context, key string) (logs []string, delim string, err error)
}
