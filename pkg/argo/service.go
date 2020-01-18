package argo

import (
	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
)

// Service is the layer of use case, it encapsulates and implements all of the use cases of the system.
type Service struct {
	repo Repository
}

// NewService create a new service.
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// Get get the workflow by the key, the format is "namespace/key", and if doesn't exist it return nil.
func (s *Service) Get(key string) *wf.Workflow {
	return s.repo.Get(key)
}

// Search get workflows which are matched with pattern.
func (s *Service) Search(pattern string) []*wf.Workflow {
	return s.repo.Search(pattern)
}

// Delete delete the workflow by the key.
func (s *Service) Delete(key string) error {
	return s.repo.Delete(key)
}

// Logs return the slice of log which type is string and the identifier and message are separated by delimeter.
func (s *Service) Logs(key string) (logs []string, delim string, err error) {
	return s.repo.Logs(key)
}
