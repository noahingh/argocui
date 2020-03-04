package resource

// RepoType is the type of repository. It gurantee the type of value.
type RepoType string

const (
	// Mock is the type of mock repository.
	Mock RepoType = "mock"
)

// UseCase is the interface of use case.
type UseCase interface {
	GetRepoType() RepoType
	Get(key string) interface{}
	Search(pattern string) []interface{}
	Delete(key string) error
}

// Repo is the interface of repository.
type Repo interface {
	Get(key string) interface{}
	Search(pattern string) []interface{}
	Delete(key string) error
}
