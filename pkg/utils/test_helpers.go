package utils

import (
	"os"
	"path/filepath"
	
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"k8s.io/apimachinery/pkg/runtime"
)

// GetTestClient returns a client for testing purposes
// It first tries to get a real client from kubeconfig
// If that fails, it returns a fake client
func GetTestClient(scheme *runtime.Scheme) (client.Client, error) {
	// Try to get a real client first
	config, err := GetTestConfig()
	if err == nil {
		return client.New(config, client.Options{Scheme: scheme})
	}
	
	// If we can't get a real client, return a fake one
	return fake.NewClientBuilder().WithScheme(scheme).Build(), nil
}

// GetTestConfig returns a rest.Config for testing purposes
func GetTestConfig() (*rest.Config, error) {
	// Try to use the in-cluster config first
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}
	
	// If KUBECONFIG is set, use that
	if kubeconfigPath := os.Getenv("KUBECONFIG"); kubeconfigPath != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}
	
	// Try the default kubeconfig path
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	kubeconfigPath := filepath.Join(home, ".kube", "config")
	if _, err := os.Stat(kubeconfigPath); err == nil {
		return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}
	
	// If we can't find a kubeconfig, check if we're in test mode
	// and return a fake config
	if os.Getenv("TEST_MODE") == "true" {
		// Return a minimal fake config
		return &rest.Config{
			Host: "http://localhost:8080",
		}, nil
	}
	
	return nil, err
}
