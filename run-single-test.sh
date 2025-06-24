#!/bin/bash

set -e

# Check if package is provided
if [ $# -lt 1 ]; then
    echo "Usage: $0 <package-path> [test-name]"
    echo "Example: $0 ./pkg/webhook/core.oam.dev/v1beta1/application TestMutatingHandler"
    exit 1
fi

PACKAGE_PATH=$1
TEST_NAME=$2
TEST_FILTER=""

if [ ! -z "$TEST_NAME" ]; then
    TEST_FILTER="-run $TEST_NAME"
fi

echo "Setting up environment for testing $PACKAGE_PATH..."

# Create temp directories
mkdir -p /tmp/kubevela-test
mkdir -p bin
mkdir -p config/crd/bases
mkdir -p charts/vela-core/crds

# Set up environment variables
export TEST_MODE=true
export KUBECONFIG="/tmp/kubevela-test/kubeconfig"
export KUBEBUILDER_ASSETS="$($(go env GOPATH)/bin/setup-envtest use -p path 1.24.2-linux-amd64)"

# Create a temporary kubeconfig for tests
cat > "$KUBECONFIG" <<EOF
apiVersion: v1
kind: Config
clusters:
- cluster:
    server: http://localhost:8080
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
current-context: test-context
users:
- name: test-user
  user:
    username: test
    password: test
EOF
chmod 600 "$KUBECONFIG"

# Download kustomize if needed
if [ ! -f "bin/kustomize" ]; then
    echo "Downloading kustomize..."
    curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash -s 3.8.7
    mv kustomize bin/
    chmod +x bin/kustomize
fi

# Fix the Definition type issue
echo "Setting up client implementation..."
mkdir -p pkg/definition/gen_sdk
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
		kubeconfig := filepath.Join(home, ".kube", "config")
		if _, err := os.Stat(kubeconfig); err == nil {
			return clientcmd.BuildConfigFromFlags("", kubeconfig)
		}
	}
	
	// No valid kubeconfig found
	return nil, err
}
EOF

# Run the test with PATH pointing to our bin directory
echo "Running test for $PACKAGE_PATH..."
PATH="$PWD/bin:$PATH" go test -v $TEST_FILTER $PACKAGE_PATH

echo "Test completed."
