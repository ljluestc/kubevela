package definition

import (
	"errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/common"
	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
)

var (
	// ErrNotDefinition is returned when the object is not a definition
	ErrNotDefinition = errors.New("not a definition object")
)

// GetGVK returns the GroupVersionKind for the definition
func (def *Definition) GetGVK() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   def.GroupVersionKind().Group,
		Version: def.GroupVersionKind().Version,
		Kind:    def.GroupVersionKind().Kind,
	}
}

// IsWorkloadDefinition checks if the definition is a WorkloadDefinition
func (def *Definition) IsWorkloadDefinition() bool {
	return def.GetKind() == "WorkloadDefinition" || def.GetKind() == "ComponentDefinition"
}

// IsTraitDefinition checks if the definition is a TraitDefinition
func (def *Definition) IsTraitDefinition() bool {
	return def.GetKind() == "TraitDefinition"
}

// IsPolicyDefinition checks if the definition is a PolicyDefinition
func (def *Definition) IsPolicyDefinition() bool {
	return def.GetKind() == "PolicyDefinition"
}

// IsWorkflowStepDefinition checks if the definition is a WorkflowStepDefinition
func (def *Definition) IsWorkflowStepDefinition() bool {
	return def.GetKind() == "WorkflowStepDefinition"
}

// ConvertWorkloadGVK2Definition converts a WorkloadDefinition to a Definition
func ConvertWorkloadGVK2Definition(obj runtime.Object) (*Definition, error) {
	// This is a placeholder implementation
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	return &Definition{
		Unstructured: unstructured.Unstructured{
			Object: unstructuredObj,
		},
	}, nil
}

// GetCUETemplate returns the CUE template from the definition's schematic
func (def *Definition) GetCUETemplate() (string, error) {
	// This is a placeholder implementation
	if def.IsWorkloadDefinition() {
		return "", nil
	}
	if def.IsTraitDefinition() {
		return "", nil
	}
	if def.IsPolicyDefinition() {
		return "", nil
	}
	if def.IsWorkflowStepDefinition() {
		return "", nil
	}
	return "", fmt.Errorf("unknown definition type: %s", def.GetKind())
}

// ExtractTypeDescriptor extracts the WorkloadTypeDescriptor from the definition
func (def *Definition) ExtractTypeDescriptor() (common.WorkloadTypeDescriptor, error) {
	// This is a placeholder implementation
	return common.WorkloadTypeDescriptor{}, nil
}

// ExtractDefinitionRef extracts the DefinitionReference from the definition
func (def *Definition) ExtractDefinitionRef() (common.DefinitionReference, error) {
	// This is a placeholder implementation
	return common.DefinitionReference{}, nil
}

// ParseSchematic parses the Schematic from the definition
func (def *Definition) ParseSchematic() (*common.Schematic, error) {
	// This is a placeholder implementation
	return &common.Schematic{}, nil
}
