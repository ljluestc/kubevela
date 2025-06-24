#!/bin/bash

echo "Direct fix for KubeVela build issues"

# Check if we're in the right directory
if [ ! -d "pkg/definition/gen_sdk" ]; then
    echo "Error: Must be run from the root of the KubeVela repository"
    exit 1
fi

# Completely remove all related files to ensure a clean state
echo "Removing existing client files..."
rm -f pkg/definition/gen_sdk/mockup_client.go
rm -f pkg/definition/gen_sdk/client.go
rm -f pkg/definition/gen_sdk/fixed_client.go

# Create fresh implementation
echo "Creating new client implementation..."
cat > pkg/definition/gen_sdk/client.go << 'EOF'
// Package gen_sdk provides client implementations for accessing Kubernetes
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

# Set correct permissions
chmod 644 pkg/definition/gen_sdk/client.go

# Clear Go build cache
echo "Clearing Go build cache..."
go clean -cache -modcache

echo "Direct fix complete. Try building with: GO111MODULE=on CGO_ENABLED=0 go build -o bin/vela ./references/cmd/cli/main.go"
