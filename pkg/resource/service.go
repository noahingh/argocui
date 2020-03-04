package resource

import (
	"context"
	"fmt"
)

var (
	// ErrNotImplement is the error when the method is not implemented.
	ErrNotImplement = fmt.Errorf("it's not implemented")
)

// Service is the layer of use case, it encapsulates and implements all of the use cases of the system.
type Service struct {
	repo Repo
}

// NewService create a new service.
func NewService(r Repo) *Service {
	return &Service{
		repo: r,
	}
}

// Get get the workflow by the key, the format is "namespace/key", and if doesn't exist it return nil.
func (s *Service) Get(key string) interface{} {
	return s.repo.Get(key)
}

// Search get workflows which are matched with pattern.
func (s *Service) Search(pattern string) []interface{} {
	return s.repo.Search(pattern)
}

// Delete delete the workflow by the key.
func (s *Service) Delete(key string) error {
	return s.repo.Delete(key)
}

// Logs get the channel to recieve Logs from a Argo workflow.
func (s *Service) Logs(ctx context.Context, key string) (ch <-chan Log, err error) {
	return s.repo.Logs(ctx, key)
}
