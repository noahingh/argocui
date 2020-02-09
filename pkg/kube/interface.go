package kube

// UseCase is the use case for Kubernetes.
type UseCase interface {
	GetNamespaces() ([]string, error)
}
