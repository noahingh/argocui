package resource

import (
	"github.com/hanjunlee/argocui/pkg/runtime/mock"
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

// GetRepoType return the type of repository.
func (s *Service) GetRepoType() RepoType {
	var t RepoType

	switch s.repo.(type) {
	case *mock.Repo:
		t = Mock
	}
	return t
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
