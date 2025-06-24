package definition

import (
	"context"
	"fmt"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/common"
	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

// ConvertDefinitionRevisionToDefinition converts a definition revision to a definition
func ConvertDefinitionRevisionToDefinition(ctx context.Context, dr runtime.Object) (runtime.Object, error) {
	switch obj := dr.(type) {
	default:
		return nil, fmt.Errorf("unknown definition revision type: %T", obj)
	}
}

// ConvertWorkloadDefinitionToComponentDefinition converts a WorkloadDefinition to a ComponentDefinition
func ConvertWorkloadDefinitionToComponentDefinition(wlDef *v1beta1.WorkloadDefinition) *v1beta1.ComponentDefinition {
	compDef := &v1beta1.ComponentDefinition{
		ObjectMeta: wlDef.ObjectMeta,
		Spec: v1beta1.ComponentDefinitionSpec{
			Workload: common.WorkloadTypeDescriptor{
				Type: wlDef.Spec.Reference.Name, // Use Reference.Name as the workload type
			},
		},
	}
	return compDef
}

// ConvertTemplateJSON2Object converts template JSON to an object
func ConvertTemplateJSON2Object(templateJSON []byte) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}
	err := yaml.Unmarshal(templateJSON, obj)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal template json: %w", err)
	}
	return obj, nil
}

// GetTypeDefReference gets type definition reference
func GetTypeDefReference(obj runtime.Object) (*common.DefinitionReference, error) {
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	u := unstructured.Unstructured{Object: unstructuredObj}

	kind := u.GetKind()
	if kind != "WorkloadDefinition" && kind != "TraitDefinition" &&
		kind != "ComponentDefinition" && kind != "PolicyDefinition" &&
		kind != "WorkflowStepDefinition" {
		return nil, fmt.Errorf("object is not a valid definition type: %s", kind)
	}

	spec, found, err := unstructured.NestedMap(u.Object, "spec")
	if err != nil || !found {
		return nil, fmt.Errorf("spec not found in definition")
	}

	ref, found, err := unstructured.NestedMap(spec, "reference")
	if err != nil || !found {
		return nil, fmt.Errorf("reference not found in spec")
	}

	name, found, err := unstructured.NestedString(ref, "name")
	if err != nil || !found {
		return nil, fmt.Errorf("name not found in reference")
	}

	return &common.DefinitionReference{
		Name: name,
	}, nil
}

// ConvertWorkloadDefinition converts a WorkloadDefinition to v1beta1.WorkloadDefinition
func ConvertWorkloadDefinition(def *Definition) (*v1beta1.WorkloadDefinition, error) {
	if !def.IsWorkloadDefinition() {
		return nil, fmt.Errorf("definition is not a WorkloadDefinition")
	}

	wd := &v1beta1.WorkloadDefinition{}
	wd.SetName(def.GetName())
	wd.SetNamespace(def.GetNamespace())
	wd.SetLabels(def.GetLabels())
	wd.SetAnnotations(def.GetAnnotations())

	spec, found, err := unstructured.NestedMap(def.Object, "spec")
	if err != nil || !found {
		return nil, fmt.Errorf("spec not found in WorkloadDefinition")
	}

	ref, found, err := unstructured.NestedMap(spec, "reference")
	if err == nil && found {
		name, found, err := unstructured.NestedString(ref, "name")
		if err == nil && found {
			wd.Spec.Reference = common.DefinitionReference{Name: name}
		}
	}

	schematic := extractSchematic(spec)
	if schematic != nil {
		wd.Spec.Schematic = schematic
	}

	return wd, nil
}

// extractSchematic extracts the schematic from the spec
func extractSchematic(spec map[string]interface{}) *common.Schematic {
	schematicData, found, err := unstructured.NestedMap(spec, "schematic")
	if err != nil || !found {
		return nil
	}

	schematic := &common.Schematic{}

	cueData, found, err := unstructured.NestedMap(schematicData, "cue")
	if err == nil && found {
		template, found, err := unstructured.NestedString(cueData, "template")
		if err == nil && found {
			schematic.CUE = &common.CUE{
				Template: template,
			}
		}
	}

	return schematic
}
