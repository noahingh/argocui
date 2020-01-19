package argo

import (
	"log"
	"os"

	"github.com/argoproj/argo/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// GetClientConfig return the client config
func GetClientConfig() clientcmd.ClientConfig {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath()},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{}})
}

// GetCurrentContext return the raw config
func GetCurrentContext() (string, error) {
	cc := GetClientConfig()
	rc, err := cc.RawConfig()
	if err != nil {
		return "", err
	}

	return rc.CurrentContext, nil
}

// GetNamespace return the namespace
func GetNamespace() (string, error) {
	cc := GetClientConfig()
	ns, _, err := cc.Namespace()
	return ns, err
}

// GetClientset return the restful API client for Argo workflow.
func GetClientset() versioned.Interface {
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath())
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := versioned.NewForConfig(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	return clientset
}

// GetKubeClientset return the restful API for Kubernetes resource.
func GetKubeClientset() kubernetes.Interface {
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath())
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	return clientset
}

func kubeConfigPath() string {
	p := os.Getenv("KUBECONFIG")
	if p == "" {
		p = os.Getenv("HOME") + "/.kube/config"
	}
	return p
}
