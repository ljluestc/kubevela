#!/bin/bash

set -e

# Display usage if no arguments provided
if [ $# -lt 1 ]; then
    echo "Usage: $0 <package-path> [test-name] [timeout]"
    echo ""
    echo "Examples:"
    echo "  $0 ./pkg/webhook/core.oam.dev/v1beta1/application           # Run all tests in package"
    echo "  $0 ./pkg/webhook/core.oam.dev/v1beta1/application TestXYZ   # Run specific test"
    echo "  $0 ./pkg/definition -timeout 5m                             # With custom timeout"
    exit 1
fi

PACKAGE_PATH=$1
TEST_NAME=$2
TIMEOUT=${3:-"2m"}

# Set up the test environment
source ./test-setup.sh

echo "================================================================"
echo "Running tests for: $PACKAGE_PATH"
if [ ! -z "$TEST_NAME" ]; then
    echo "Test name: $TEST_NAME"
    TEST_FILTER="-run $TEST_NAME"
fi
echo "Timeout: $TIMEOUT"
echo "================================================================"

# Run the test with environment variables set
PATH="$PWD/bin:$PATH" \
KUBEBUILDER_ASSETS="$($(go env GOPATH)/bin/setup-envtest use -p path 1.24.2)" \
TEST_MODE=true \
go test -v -timeout $TIMEOUT $TEST_FILTER $PACKAGE_PATH

echo "Test completed."
