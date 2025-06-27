package testing

import (
	"os"
	"path/filepath"
	"strings"
	stdtesting "testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

// TestEnv encapsulates the test environment setup
type TestEnv struct {
	Environment *envtest.Environment
	Config      *rest.Config
	Client      client.Client
	CRDPaths    []string
}

// NewTestEnv creates a new test environment with the specified CRD paths
func NewTestEnv(t *stdtesting.T, crdPaths []string) *TestEnv {
	require := require.New(t)

	env := &TestEnv{
		CRDPaths: crdPaths,
	}

	env.Environment = &envtest.Environment{
		CRDDirectoryPaths: crdPaths,
	}

	// Set up environment variables
	os.Setenv("TEST_MODE", "true")

	// Create a temporary KUBECONFIG
	tmpDir, err := os.MkdirTemp("", "test-kubeconfig")
	require.NoError(err, "should create temp directory")

	kubeconfigPath := filepath.Join(tmpDir, "kubeconfig")
	os.Setenv("KUBECONFIG", kubeconfigPath)

	// Start the test environment
	config, err := env.Environment.Start()
	require.NoError(err, "should start test environment")
	env.Config = config

	// Create test client
	client, err := client.New(config, client.Options{})
	require.NoError(err, "should create client")
	env.Client = client

	return env
}

// Stop stops the test environment and cleans up resources
func (e *TestEnv) Stop(t *stdtesting.T) {
	if e.Environment != nil {
		err := e.Environment.Stop()
		if err != nil {
			t.Errorf("failed to stop test environment: %v", err)
		}
	}

	// Clean up temporary KUBECONFIG
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig != "" {
		if strings.HasPrefix(kubeconfig, os.TempDir()) {
			err := os.Remove(kubeconfig)
			if err != nil {
				t.Logf("Error removing temp kubeconfig: %v", err)
			}
		}
	}
}

// DefaultCRDPaths returns the default CRD paths used in the project
func DefaultCRDPaths() []string {
	return []string{
		"../../config/crd/bases",
		"../../charts/vela-core/crds",
	}
}
