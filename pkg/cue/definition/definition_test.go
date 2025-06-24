package definition


func TestDefinitionWithTestObjects(t *testing.T) {
	// Create test objects
	compDef := CreateTestComponentDefinition("test-comp", "default", "template: {}")
	traitDef := CreateTestTraitDefinition("test-trait", "default", "patch: {}")
	
	// Create test client with objects
	def := SetupTestDefinitionClient(t, compDef, traitDef)
	
	// Test component definition template
	ctx := context.Background()
	compTemplate, err := def.GetDefinitionTemplate(ctx, "component", "test-comp", "default")
	assert.NoError(t, err)
	assert.Equal(t, "template: {}", compTemplate)
	
	// Test trait definition template
	traitTemplate, err := def.GetDefinitionTemplate(ctx, "trait", "test-trait", "default")
	assert.NoError(t, err)
	assert.Equal(t, "patch: {}", traitTemplate)
	
	// Test invalid definition type
	_, err = def.GetDefinitionTemplate(ctx, "invalid", "test", "default")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported definition type")
	
	// Test non-existent definition
	_, err = def.GetDefinitionTemplate(ctx, "component", "non-existent", "default")
	assert.Error(t, err)
}

func TestCreateTestDefinitions(t *testing.T) {
	// Test CreateTestComponentDefinition
	compDef := CreateTestComponentDefinition("test-comp", "default", "template: {}")
	assert.Equal(t, "test-comp", compDef.Name)
	assert.Equal(t, "default", compDef.Namespace)
	assert.Equal(t, "test-workload", compDef.Spec.Workload.Type)
	assert.Equal(t, "template: {}", compDef.Spec.Schematic.CUE.Template)
	
	// Test CreateTestTraitDefinition
	traitDef := CreateTestTraitDefinition("test-trait", "default", "patch: {}")
	assert.Equal(t, "test-trait", traitDef.Name)
	assert.Equal(t, "default", traitDef.Namespace)
	assert.Equal(t, "test-trait-ref", traitDef.Spec.Reference.Name)
	assert.Equal(t, "patch: {}", traitDef.Spec.Schematic.CUE.Template)
}

func TestNewDefinition(t *testing.T) {
	scheme := runtime.NewScheme()
	assert.NoError(t, v1beta1.AddToScheme(scheme))
	
	compDef := &v1beta1.ComponentDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-component",
			Namespace: "default",
		},
	}
	
	def := SetupTestDefinitionClient(t, compDef)
	assert.NotNil(t, def)
}
import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oam-dev/kubevela/pkg/oam/util"
)

func TestGetDefinitionTemplate(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	// Create test ComponentDefinition
	compDef := CreateTestComponentDefinition("test-comp", "default", `
parameter: {
	replicas: *1 | int
}
output: {
	apiVersion: "apps/v1"
	kind: "Deployment"
	spec: {
		replicas: parameter.replicas
	}
}
`)
// Package definition provides tests for definition utilities
package definition

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/common"
	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
)

func TestGetUserAnnotationsAndLabels(t *testing.T) {
	def := &Definition{
		Unstructured: unstructured.Unstructured{},
	}
	
	// Set annotations
	def.SetAnnotations(map[string]string{
		UserPrefix + "annotation": "test-value",
		"other-annotation":        "other-value",
	})
	
	// Set labels
	def.SetLabels(map[string]string{
		UserPrefix + "label": "test-label",
		"other-label":        "other-label",
	})
	
	// Test GetUserAnnotations
	userAnnotations := def.GetUserAnnotations()
	assert.Equal(t, 1, len(userAnnotations))
	assert.Equal(t, "test-value", userAnnotations[UserPrefix+"annotation"])
	
	// Test GetUserLabels
	userLabels := def.GetUserLabels()
	assert.Equal(t, 1, len(userLabels))
	assert.Equal(t, "test-label", userLabels[UserPrefix+"label"])
}

func TestDefinitionWithTestObjects(t *testing.T) {
	// Create test objects
	compDef := CreateTestComponentDefinition("test-comp", "default", "template: {}")
	traitDef := CreateTestTraitDefinition("test-trait", "default", "patch: {}")
	
	// Create test client with objects
	def := SetupTestDefinitionClient(t, compDef, traitDef)
	
	// Test component definition template
	ctx := context.Background()
	compTemplate, err := def.GetDefinitionTemplate(ctx, "component", "test-comp", "default")
	assert.NoError(t, err)
	assert.Equal(t, "template: {}", compTemplate)
	
	// Test trait definition template
	traitTemplate, err := def.GetDefinitionTemplate(ctx, "trait", "test-trait", "default")
	assert.NoError(t, err)
	assert.Equal(t, "patch: {}", traitTemplate)
	
	// Test invalid definition type
	_, err = def.GetDefinitionTemplate(ctx, "invalid", "test", "default")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported definition type")
	
	// Test non-existent definition
	_, err = def.GetDefinitionTemplate(ctx, "component", "non-existent", "default")
	assert.Error(t, err)
}

func TestCreateTestDefinitions(t *testing.T) {
	// Test CreateTestComponentDefinition
	compDef := CreateTestComponentDefinition("test-comp", "default", "template: {}")
	assert.Equal(t, "test-comp", compDef.Name)
	assert.Equal(t, "default", compDef.Namespace)
	assert.Equal(t, "test-workload", compDef.Spec.Workload.Type)
	assert.Equal(t, "template: {}", compDef.Spec.Schematic.CUE.Template)
	
	// Test CreateTestTraitDefinition
	traitDef := CreateTestTraitDefinition("test-trait", "default", "patch: {}")
	assert.Equal(t, "test-trait", traitDef.Name)
	assert.Equal(t, "default", traitDef.Namespace)
	assert.Equal(t, "test-trait-ref", traitDef.Spec.Reference.Name)
	assert.Equal(t, "patch: {}", traitDef.Spec.Schematic.CUE.Template)
}

func TestConvertWorkloadDefinitionToComponentDefinition(t *testing.T) {
	wlDef := &v1beta1.WorkloadDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-workload",
			Namespace: "default",
		},
		Spec: v1beta1.WorkloadDefinitionSpec{
			Reference: common.Reference{
				Name: "test-type",
			},
		},
	}
	
	compDef := ConvertWorkloadDefinitionToComponentDefinition(wlDef)
	
	assert.Equal(t, "test-workload", compDef.Name)
	assert.Equal(t, "default", compDef.Namespace)
	assert.Equal(t, "test-type", compDef.Spec.Workload.Type)
}

func TestIsWorkloadDefinition(t *testing.T) {
	def := &Definition{
		Unstructured: unstructured.Unstructured{},
	}
	def.SetKind("WorkloadDefinition")
	assert.True(t, def.IsWorkloadDefinition())
	
	def.SetKind("ComponentDefinition")
	assert.False(t, def.IsWorkloadDefinition())
}
	// Create test TraitDefinition
	traitDef := CreateTestTraitDefinition("test-trait", "default", `
parameter: {
	annotations: [string]: string
}
patch: {
	metadata: {
		annotations: parameter.annotations
	}
}
`)

	// Setup test client with definitions
	def := SetupTestDefinitionClient(t, compDef, traitDef)

	// Test getting ComponentDefinition template
	template, err := def.GetDefinitionTemplate(ctx, util.DefinitionTypeComponent, "test-comp", "default")
	assert.NoError(err)
	assert.Contains(template, "parameter")
	assert.Contains(template, "replicas")

	// Test getting TraitDefinition template
	template, err = def.GetDefinitionTemplate(ctx, util.DefinitionTypeTrait, "test-trait", "default")
	assert.NoError(err)
	assert.Contains(template, "parameter")
	assert.Contains(template, "annotations")

	// Test getting nonexistent definition
	_, err = def.GetDefinitionTemplate(ctx, util.DefinitionTypeComponent, "nonexistent", "default")
	assert.Error(err)
}

func TestGetDefinitionCapabilities(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
package definition

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
	"github.com/oam-dev/kubevela/pkg/utils/common"
)

func TestDefinitionWithClient(t *testing.T) {
	// Create a test ComponentDefinition
	compDef := CreateTestComponentDefinition("test-component", "default", `
output: {
    apiVersion: "apps/v1"
    kind:       "Deployment"
    spec: {
        selector: matchLabels: {
            "app.oam.dev/component": context.name
        }
        template: {
            metadata: labels: {
                "app.oam.dev/component": context.name
            }
            spec: {
                containers: [{
                    name:  context.name
                    image: parameter.image
                    if parameter["cmd"] != _|_ {
                        command: parameter.cmd
                    }
                }]
            }
        }
    }
}
parameter: {
    image: string
    cmd?: [...string]
}
`)

	// Create a test TraitDefinition
	traitDef := CreateTestTraitDefinition("test-trait", "default", `
patch: {
   spec: replicas: parameter.replicas
}
parameter: {
    replicas: *1 | int
}
`)

	// Set up the test client with the definitions
	def := SetupTestDefinitionClient(t, compDef, traitDef)

	// Test getting templates
	ctx := context.Background()
	
	// Test component template
	compTemplate, err := def.GetDefinitionTemplate(ctx, "component", "test-component", "default")
	assert.NoError(t, err)
	assert.Contains(t, compTemplate, "parameter: {")
	
	// Test trait template
	traitTemplate, err := def.GetDefinitionTemplate(ctx, "trait", "test-trait", "default")
	assert.NoError(t, err)
	assert.Contains(t, traitTemplate, "replicas: *1 | int")
	
	// Test error case - non-existent definition
	_, err = def.GetDefinitionTemplate(ctx, "component", "non-existent", "default")
	assert.Error(t, err)
}

func TestDefinitionBasicFunctionsWithStructs(t *testing.T) {
	c := fake.NewClientBuilder().WithScheme(common.Scheme).Build()
	def := &Definition{
		Unstructured: unstructured.Unstructured{},
		client:       c,
	}
	def.SetAnnotations(map[string]string{
		UserPrefix + "annotation": "annotation",
		"other":                   "other",
	})
	def.SetLabels(map[string]string{
		UserPrefix + "label": "label",
		"other":              "other",
	})
	assert.Equal(t, map[string]string{UserPrefix + "annotation": "annotation"}, def.GetUserAnnotations())
	assert.Equal(t, map[string]string{UserPrefix + "label": "label"}, def.GetUserLabels())
}

func TestNewDefinitionWithClient(t *testing.T) {
	scheme := runtime.NewScheme()
	assert.NoError(t, v1beta1.AddToScheme(scheme))
	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	def := NewDefinition(client)
	assert.NotNil(t, def)
	assert.NotNil(t, def.client)
}
	// Create test definitions
	compDef := CreateTestComponentDefinition("test-comp", "default", `output: {}`)
	traitDef := CreateTestTraitDefinition("test-trait", "default", `patch: {}`)

	// Setup test client with definitions
	def := SetupTestDefinitionClient(t, compDef, traitDef)

	// Test getting ComponentDefinition capability
	capabilities, err := def.GetCapabilitiesFromDefinition(ctx, util.DefinitionTypeComponent, "test-comp", "default")
	assert.NoError(err)
	assert.Equal(1, len(capabilities))
	assert.Equal("test-comp", capabilities[0].Name)
	assert.Equal("Component", capabilities[0].Type)

	// Test getting TraitDefinition capability
	capabilities, err = def.GetCapabilitiesFromDefinition(ctx, util.DefinitionTypeTrait, "test-trait", "default")
	assert.NoError(err)
	assert.Equal(1, len(capabilities))
	assert.Equal("test-trait", capabilities[0].Name)
	assert.Equal("Trait", capabilities[0].Type)

	// Test getting nonexistent definition
	_, err = def.GetCapabilitiesFromDefinition(ctx, util.DefinitionTypeComponent, "nonexistent", "default")
	assert.Error(err)
}
