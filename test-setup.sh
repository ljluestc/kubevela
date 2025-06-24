#!/bin/sh
# Script to set up environment for running tests

# Ensure we have a valid kubeconfig for testing
if [ -z "$KUBECONFIG" ] && [ -z "$KUBERNETES_MASTER" ]; then
    echo "No Kubernetes configuration found. Creating a temporary one for testing..."
    
    # Create temp directory for test artifacts if it doesn't exist
    TEST_TMP_DIR="/tmp/kubevela-test"
    mkdir -p "$TEST_TMP_DIR"
    
    # Set up a temporary kubeconfig file
    export KUBECONFIG="$TEST_TMP_DIR/kubeconfig"
    
    # Create a minimal kubeconfig file for tests
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
    
    echo "Temporary kubeconfig created at $KUBECONFIG"
fi

# Set other environment variables needed for tests
export TEST_MODE=true

echo "Test environment setup complete. You can now run tests with:"
echo "go test ./pkg/... -v"
