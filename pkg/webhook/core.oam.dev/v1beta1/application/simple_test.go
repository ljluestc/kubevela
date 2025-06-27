package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
)

// TestSimpleApplication is a minimal test to verify the environment is set up correctly
func TestSimpleApplication(t *testing.T) {
	// Create a test scheme
	scheme := runtime.NewScheme()
	assert.NoError(t, v1beta1.AddToScheme(scheme))

	// Create a simple application
	app := &v1beta1.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-app",
			Namespace: "default",
		},
		Spec: v1beta1.ApplicationSpec{
			Components: []v1beta1.ApplicationComponent{
				{
					Name: "test-component",
					Type: "webservice",
				},
			},
		},
	}

	// Create a fake client
	client := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(app).
		Build()

	// Verify we can get the application
	result := &v1beta1.Application{}
	err := client.Get(ctx, client.ObjectKey{
		Name:      "test-app",
		Namespace: "default",
	}, result)

	// Check the results
	assert.NoError(t, err)
	assert.Equal(t, "test-app", result.Name)
	assert.Equal(t, 1, len(result.Spec.Components))
	assert.Equal(t, "test-component", result.Spec.Components[0].Name)
}
