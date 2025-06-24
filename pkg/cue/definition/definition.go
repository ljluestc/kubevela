package definition

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
)

// UserPrefix is the prefix for user-provided annotations and labels
const UserPrefix = "custom.definition.oam.dev/"

// Definition is a wrapper for working with definition resources
type Definition struct {
	unstructured.Unstructured
	client client.Client
}

// NewDefinition creates a new Definition with a client
func NewDefinition(c client.Client) *Definition {
	return &Definition{
		client: c,
	}
}

// NewDefinitionFromObject creates a new Definition from an unstructured object
func NewDefinitionFromObject(obj runtime.Object) (*Definition, error) {
	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	return &Definition{Unstructured: unstructured.Unstructured{Object: u}}, nil
}

// GetUserAnnotations gets annotations with UserPrefix
func (def *Definition) GetUserAnnotations() map[string]string {
	return filterMapWithPrefix(def.GetAnnotations(), UserPrefix)
}

// GetUserLabels gets labels with UserPrefix
func (def *Definition) GetUserLabels() map[string]string {
	return filterMapWithPrefix(def.GetLabels(), UserPrefix)
}

// filterMapWithPrefix returns a map with only keys that have the specified prefix
func filterMapWithPrefix(m map[string]string, prefix string) map[string]string {
	if m == nil {
		return nil
	}
	filtered := map[string]string{}
	for k, v := range m {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			filtered[k] = v
		}
	}
	return filtered
}

// GetDefinitionTemplate retrieves the CUE template from the specified definition
func (def *Definition) GetDefinitionTemplate(ctx context.Context, defType, name, namespace string) (string, error) {
	var template string
	var err error

	switch defType {
	case "component":
		template, err = def.getComponentDefinitionTemplate(ctx, name, namespace)
	case "trait":
		template, err = def.getTraitDefinitionTemplate(ctx, name, namespace)
	default:
		return "", fmt.Errorf("unsupported definition type: %s", defType)
	}

	if err != nil {
		return "", err
	}

	return template, nil
}

// getComponentDefinitionTemplate retrieves the CUE template from a ComponentDefinition
func (def *Definition) getComponentDefinitionTemplate(ctx context.Context, name, namespace string) (string, error) {
	compDef := &v1beta1.ComponentDefinition{}
	if err := def.client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, compDef); err != nil {
		return "", fmt.Errorf("failed to get ComponentDefinition %s/%s: %w", namespace, name, err)
	}

	if compDef.Spec.Schematic == nil || compDef.Spec.Schematic.CUE == nil {
		return "", fmt.Errorf("ComponentDefinition %s/%s does not have a CUE template", namespace, name)
	}

	return compDef.Spec.Schematic.CUE.Template, nil
}

// getTraitDefinitionTemplate retrieves the CUE template from a TraitDefinition
func (def *Definition) getTraitDefinitionTemplate(ctx context.Context, name, namespace string) (string, error) {
	traitDef := &v1beta1.TraitDefinition{}
	if err := def.client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, traitDef); err != nil {
		return "", fmt.Errorf("failed to get TraitDefinition %s/%s: %w", namespace, name, err)
	}

	if traitDef.Spec.Schematic == nil || traitDef.Spec.Schematic.CUE == nil {
		return "", fmt.Errorf("TraitDefinition %s/%s does not have a CUE template", namespace, name)
	}

	return traitDef.Spec.Schematic.CUE.Template, nil
}

// ToApplication converts the definition to an Application
func (def *Definition) ToApplication() (*v1beta1.Application, error) {
	return nil, nil // Placeholder for implementation
}

// GetComponent returns the component from a definition
func (def *Definition) GetComponent() (interface{}, error) {
	return nil, nil // Placeholder for implementation
}
