#!/bin/bash

echo "Fixing client issues in KubeVela..."

# Remove the problematic file
echo "Removing problematic file..."
rm -f pkg/definition/gen_sdk/mockup_client.go

# Create a new file with the correct implementation
echo "Creating new fixed client file..."
cat > pkg/definition/gen_sdk/client.go << 'EOF'
package gen_sdk

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// GetMockClient returns a fake client for testing purposes
func GetMockClient(s *runtime.Scheme) (client.Client, error) {
	return fake.NewClientBuilder().WithScheme(s).Build(), nil
}

// GetClient attempts to get a real Kubernetes client or falls back to a mock client
func GetClient() (client.Client, error) {
	// Try to get a real client first
	config, err := getKubeConfig()
	if err != nil {
		// If we can't get a real client and we're in test mode, use a mock
		if os.Getenv("TEST_MODE") == "true" {
			// Using scheme variable directly without function call
			k8sScheme := scheme.Scheme
			return GetMockClient(k8sScheme)
		}
		return nil, err
	}
	
	// Using scheme variable directly without function call
	k8sScheme := scheme.Scheme
	
	return client.New(config, client.Options{Scheme: k8sScheme})
}

// getKubeConfig tries to find and load a kubeconfig file
func getKubeConfig() (*rest.Config, error) {
	// Try in-cluster config first
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}
	
	// Try KUBECONFIG env var
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig != "" {
		if _, err := os.Stat(kubeconfig); err == nil {
			return clientcmd.BuildConfigFromFlags("", kubeconfig)
		}
	}
	
	// Try default location
	home, err := os.UserHomeDir()
	if err == nil {
		kubeconfig = filepath.Join(home, ".kube", "config")
		if _, err := os.Stat(kubeconfig); err == nil {
			return clientcmd.BuildConfigFromFlags("", kubeconfig)
		}
	}
	
	// No valid kubeconfig found
	return nil, err
}
EOF

# Make sure the new file has the correct permissions
chmod 644 pkg/definition/gen_sdk/client.go

echo "Fix complete. Now try running your build command again."
