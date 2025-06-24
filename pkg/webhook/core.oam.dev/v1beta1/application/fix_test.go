// Package application contains webhook handlers for Application resources
package application

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// SimpleTest is a minimal test to verify test environment setup
func SimpleTest(t *testing.T) {
	assert.True(t, true, "True should be true")
	fmt.Println("Simple test passed!")
}
