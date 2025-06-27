package testing

import (
	"context"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

// GetMockClient returns a fake client for testing purposes
func GetMockClient(scheme *runtime.Scheme, initObjs ...runtime.Object) client.Client {
	builder := fake.NewClientBuilder()
	if scheme != nil {
		builder = builder.WithScheme(scheme)
	} else {
		builder = builder.WithScheme(scheme.Scheme)
	}
	
	if len(initObjs) > 0 {
		builder = builder.WithRuntimeObjects(initObjs...)
	}
	
	return builder.Build()
}

// SetupEnvTest creates a test environment with the specified CRDs
func SetupEnvTest(crdPaths []string) (*envtest.Environment, error) {
	testEnv := &envtest.Environment{
		CRDDirectoryPaths: crdPaths,
	}
	
	// Set TEST_MODE to ensure other components know we're in test mode
	os.Setenv("TEST_MODE", "true")
	
	// Create a temporary KUBECONFIG if not set
	if os.Getenv("KUBECONFIG") == "" {
		tmpDir, err := os.MkdirTemp("", "kubeconfig")
		if err != nil {
			return nil, err
		}
		
		kubeconfig := filepath.Join(tmpDir, "kubeconfig")
		os.Setenv("KUBECONFIG", kubeconfig)
	}
	
	cfg, err := testEnv.Start()
	if err != nil {
		return nil, err
	}
	
	return testEnv, nil
}

// GetClientFromEnv creates a real client from the test environment
func GetClientFromEnv(env *envtest.Environment, scheme *runtime.Scheme) (client.Client, error) {
	cfg, err := env.Config()
	if err != nil {
		return nil, err
	}
	
	if scheme == nil {
		scheme = runtime.NewScheme()
		_ = scheme.Scheme.AddToScheme(scheme)
	}
	
	return client.New(cfg, client.Options{Scheme: scheme})
}

// CreateObject is a helper to create an object and wait for it to exist
func CreateObject(ctx context.Context, c client.Client, obj client.Object) error {
	err := c.Create(ctx, obj)
	if err != nil {
		return err
	}
	
	// Get the object to ensure it was created
	return c.Get(ctx, client.ObjectKeyFromObject(obj), obj)
}

// DeleteObject is a helper to delete an object
func DeleteObject(ctx context.Context, c client.Client, obj client.Object) error {
	return c.Delete(ctx, obj)
}

// GetTestConfig returns a REST config for testing
func GetTestConfig() (*rest.Config, error) {
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
	
	// For testing, return a minimal config
	if os.Getenv("TEST_MODE") == "true" {
		return &rest.Config{
			Host: "http://localhost:8080",
		}, nil
	}
	
	return nil, err
}
