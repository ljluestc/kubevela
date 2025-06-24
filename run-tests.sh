#!/bin/sh
# Convenient script to run tests with proper environment setup

# Source the test setup script to prepare the environment
. ./test-setup.sh

# Run the tests with the specified arguments or default to all packages
if [ $# -eq 0 ]; then
    echo "Running all tests..."
    go test ./pkg/... -v
else
    echo "Running specified tests..."
    go test "$@"
fi
