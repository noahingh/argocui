package kube

import (
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Service is the implement of the use case.
//
// The the clientset of Kubernetes is the repository of service util it needs a abundant features.
type Service struct {
	repo kubernetes.Interface
}

// NewService create a new service of Kubernetes.
func NewService(clientset kubernetes.Interface) *Service {
	return &Service{
		repo: clientset,
	}
}

// GetNamespaces return the namespaces of the cluster.
func (s *Service) GetNamespaces() ([]string, error) {
	var (
		namespaces = make([]string, 0)
	)

	list, err := s.repo.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, ns := range list.Items {
		namespaces = append(namespaces, ns.GetName())
	}

	return namespaces, nil
}
