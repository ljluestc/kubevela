#!/bin/bash

# This script verifies that the necessary tools are installed for testing

echo "Verifying installation..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed or not in PATH"
    exit 1
fi

echo "Go version: $(go version)"

# Check if GOPATH is set
if [ -z "$GOPATH" ]; then
    echo "GOPATH is not set"
    exit 1
fi

echo "GOPATH: $GOPATH"

# Check if setup-envtest is installed
if ! command -v setup-envtest &> /dev/null; then
    echo "Installing setup-envtest..."
    go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
else
    echo "setup-envtest is already installed"
fi

# Create bin directory if it doesn't exist
if [ ! -d "bin" ]; then
    mkdir -p bin
    echo "Created bin directory"
fi

# Install kustomize if not present
if [ ! -f "bin/kustomize" ]; then
    echo "Installing kustomize..."
    curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash -s 3.8.7
    mv kustomize bin/
    chmod +x bin/kustomize
    echo "Kustomize installed to bin/kustomize"
else
    echo "Kustomize is already installed in bin/kustomize"
fi

# Verify kustomize installation
if [ -f "bin/kustomize" ]; then
    echo "Kustomize version: $(bin/kustomize version --short)"
else
    echo "Kustomize installation failed"
    exit 1
fi

# Set up kubebuilder assets
export KUBEBUILDER_ASSETS="$(setup-envtest use -p path 1.24.2-linux-amd64)"
echo "KUBEBUILDER_ASSETS: $KUBEBUILDER_ASSETS"

# Create necessary directories
mkdir -p config/crd/bases
mkdir -p charts/vela-core/crds

echo "Installation verification complete."
echo "To run a single test package, use: ./run-single-test.sh <package-path>"
echo "Example: ./run-single-test.sh ./pkg/webhook/core.oam.dev/v1beta1/application"
