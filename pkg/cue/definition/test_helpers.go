package definition

import (
	"context"
	"testing"

	"cuelang.org/go/cue"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
)

// SetupTestDefinitionClient creates a fake client with test definitions
func SetupTestDefinitionClient(t *testing.T, defs ...runtime.Object) *Definition {
	assert := assert.New(t)

	// Create a scheme with the necessary types
	scheme := runtime.NewScheme()
	assert.NoError(v1beta1.AddToScheme(scheme))

	// Convert runtime.Object to client.Object for the client builder
	clientObjs := make([]client.Object, 0, len(defs))
	for _, obj := range defs {
		clientObjs = append(clientObjs, obj.(client.Object))
	}

	// Create a fake client with the test objects
	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(clientObjs...).
		Build()

	// Create a new Definition with the fake client
	def := NewDefinition(fakeClient)
	assert.NotNil(def)

	return def
}

// CreateTestComponentDefinition creates a test ComponentDefinition
func CreateTestComponentDefinition(name, namespace string, template string) *v1beta1.ComponentDefinition {
	return &v1beta1.ComponentDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1beta1.ComponentDefinitionSpec{
			Workload: v1beta1.WorkloadTypeDescriptor{
				Type: "test-workload",
			},
			Schematic: &v1beta1.Schematic{
				CUE: &v1beta1.CUE{
					Template: template,
				},
			},
		},
	}
}

// CreateTestTraitDefinition creates a test TraitDefinition
func CreateTestTraitDefinition(name, namespace string, template string) *v1beta1.TraitDefinition {
	return &v1beta1.TraitDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1beta1.TraitDefinitionSpec{
			Reference: v1beta1.DefinitionReference{
				Name: name + "-ref",
			},
			Schematic: &v1beta1.Schematic{
				CUE: &v1beta1.CUE{
					Template: template,
				},
			},
		},
	}
}

// AssertDefinitionTemplate tests the Definition's template parsing functionality
func AssertDefinitionTemplate(t *testing.T, def *Definition, defType, name, namespace string) {
	assert := assert.New(t)

	ctx := context.Background()
	template, err := def.GetDefinitionTemplate(ctx, defType, name, namespace)
	assert.NoError(err)
	assert.NotEmpty(template)

	// Test that the template can be compiled
	r := cue.Runtime{}
	inst, err := r.Compile("test", template)
	assert.NoError(err)
	assert.NotNil(inst)
}
