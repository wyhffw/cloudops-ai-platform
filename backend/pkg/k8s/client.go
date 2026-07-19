package k8s

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClient(kubeconfig string) (*kubernetes.Clientset, error) {
	cfg, err := restConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(cfg)
}

func restConfig(kubeconfig string) (*rest.Config, error) {
	// In-cluster first when running as a Pod.
	if cfg, err := rest.InClusterConfig(); err == nil {
		return cfg, nil
	}

	if kubeconfig == "" {
		if home, err := os.UserHomeDir(); err == nil {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}
	if kubeconfig == "" {
		return nil, fmt.Errorf("kubeconfig not found; set KUBECONFIG or run in-cluster")
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}
