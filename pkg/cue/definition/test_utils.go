package definition

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
	"github.com/oam-dev/kubevela/pkg/utils/common"
)

// CreateTestNamespace creates a test namespace
func CreateTestNamespace(name string) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

// SetupTestEnv sets up a test environment with namespaces and definitions
func SetupTestEnv(t *testing.T) *Definition {
	// Create a test namespace
	ns := CreateTestNamespace("test-ns")
	
	// Create a test ComponentDefinition
	compDef := CreateTestComponentDefinition("test-component", "test-ns", `
output: {
    apiVersion: "apps/v1"
    kind:       "Deployment"
    spec: {
        selector: matchLabels: {
            "app.oam.dev/component": context.name
        }
    }
}
parameter: {
    replicas: *1 | int
}
`)
	
	// Create a test TraitDefinition
	traitDef := CreateTestTraitDefinition("test-trait", "test-ns", `
patch: {
   spec: replicas: parameter.replicas
}
parameter: {
    replicas: *1 | int
}
`)
	
	// Set up the test client with namespace and definitions
	return SetupTestDefinitionClient(t, ns, compDef, traitDef)
}

// CreateMockClient creates a mock client with the given objects
func CreateMockClient(objs ...runtime.Object) (client.Client, error) {
	scheme := runtime.NewScheme()
	if err := v1beta1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	
	// Convert runtime.Object to client.Object
	clientObjs := make([]client.Object, 0, len(objs))
	for _, obj := range objs {
		clientObjs = append(clientObjs, obj.(client.Object))
	}
	
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(clientObjs...).Build(), nil
}

// CreateUnstructuredObj creates an unstructured object from a GVK and data
func CreateUnstructuredObj(group, version, kind, name, namespace string) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	})
	obj.SetName(name)
	obj.SetNamespace(namespace)
	return obj
}

// TestGetDefinition tests getting a definition from the client
func TestGetDefinition(t *testing.T, c client.Client, name, namespace, kind string) {
	assert := assert.New(t)
	
	ctx := context.Background()
	key := types.NamespacedName{Name: name, Namespace: namespace}
	
	switch kind {
	case "ComponentDefinition":
		obj := &v1beta1.ComponentDefinition{}
		err := c.Get(ctx, key, obj)
		assert.NoError(err)
		assert.Equal(name, obj.Name)
	case "TraitDefinition":
		obj := &v1beta1.TraitDefinition{}
		err := c.Get(ctx, key, obj)
		assert.NoError(err)
		assert.Equal(name, obj.Name)
	default:
		assert.Fail(fmt.Sprintf("Unsupported kind: %s", kind))
	}
}

// TestNewClientSchemeBuilder tests creating a client with a scheme builder
func TestNewClientSchemeBuilder(t *testing.T) {
	scheme := common.Scheme
	cm := &corev1.ConfigMap{}
	secret := &corev1.Secret{}
	defs := []client.Object{cm, secret}
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(defs...).Build()
	
	def := NewDefinition(client)
	assert.NotNil(t, def)
}
