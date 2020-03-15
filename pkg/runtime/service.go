package runtime

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
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

// Get get the object by the key, the format is "namespace/key", and if doesn't exist it return nil.
func (s *Service) Get(key string) (runtime.Object, error) {
	return s.repo.Get(key)
}

// Search get objects which are matched with pattern.
func (s *Service) Search(namespace, pattern string) []runtime.Object {
	return s.repo.Search(namespace, pattern)
}

// Delete delete the object by the key.
func (s *Service) Delete(key string) error {
	return s.repo.Delete(key)
}

// Logs follow logs from the object until the context is done.
func (s *Service) Logs(ctx context.Context, key string) (<-chan Log, error) {
	return s.repo.Logs(ctx, key)
}
