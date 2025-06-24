package definition

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
)

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
func GetTypeDefReference(obj runtime.Object) (*v1beta1.DefinitionReference, error) {
	// Convert the object to unstructured first
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	u := unstructured.Unstructured{Object: unstructuredObj}

	// Check if the object is a valid definition type
	kind := u.GetKind()
	if kind != "WorkloadDefinition" && kind != "TraitDefinition" && 
	   kind != "ComponentDefinition" && kind != "PolicyDefinition" && 
	   kind != "WorkflowStepDefinition" {
		return nil, fmt.Errorf("object is not a valid definition type: %s", kind)
	}

	// Extract the definitionRef field from spec
	spec, found, err := unstructured.NestedMap(u.Object, "spec")
	if err != nil || !found {
		return nil, fmt.Errorf("spec not found in definition")
	}

	defRef, found, err := unstructured.NestedMap(spec, "definitionRef")
	if err != nil || !found {
		return nil, fmt.Errorf("definitionRef not found in spec")
	}

	name, found, err := unstructured.NestedString(defRef, "name")
	if err != nil || !found {
		return nil, fmt.Errorf("name not found in definitionRef")
	}

	version, found, err := unstructured.NestedString(defRef, "version")
	if err != nil {
		return nil, fmt.Errorf("error getting version from definitionRef: %w", err)
	}
	if !found {
		version = ""
	}

	return &v1beta1.DefinitionReference{
		Name:    name,
		Version: version,
	}, nil
}

// ConvertWorkloadDefinition converts a WorkloadDefinition to v1beta1.WorkloadDefinition
func ConvertWorkloadDefinition(def *Definition) (*v1beta1.WorkloadDefinition, error) {
	if !def.IsWorkloadDefinition() {
		return nil, fmt.Errorf("definition is not a WorkloadDefinition")
	}

	// Create a new v1beta1.WorkloadDefinition
	wd := &v1beta1.WorkloadDefinition{}
	
	// Set basic metadata
	wd.SetName(def.GetName())
	wd.SetNamespace(def.GetNamespace())
	wd.SetLabels(def.GetLabels())
	wd.SetAnnotations(def.GetAnnotations())

	// Extract workloadType
	spec, found, err := unstructured.NestedMap(def.Object, "spec")
	if err != nil || !found {
		return nil, fmt.Errorf("spec not found in WorkloadDefinition")
	}

	workloadType, found, err := unstructured.NestedString(spec, "workloadType")
	if err != nil || !found {
		return nil, fmt.Errorf("workloadType not found in spec")
	}
	
	// Set the workload type descriptor
	wd.Spec.WorkloadType = workloadType

	// Extract definitionRef if it exists
	defRef, err := GetTypeDefReference(def)
	if err == nil {
		wd.Spec.DefinitionRef = *defRef
	}

	// Extract schematic
	schematic, err := extractSchematic(spec)
	if err == nil {
		wd.Spec.Schematic = *schematic
	}

	return wd, nil
}

// extractSchematic extracts the schematic from the spec
func extractSchematic(spec map[string]interface{}) (*v1beta1.Schematic, error) {
	schematicData, found, err := unstructured.NestedMap(spec, "schematic")
	if err != nil || !found {
		return nil, fmt.Errorf("schematic not found in spec")
	}

	schematic := &v1beta1.Schematic{}

	// Check for CUE schematic
	cueData, found, err := unstructured.NestedMap(schematicData, "cue")
	if err == nil && found {
		template, found, err := unstructured.NestedString(cueData, "template")
		if err != nil || !found {
			return nil, fmt.Errorf("template not found in cue schematic")
		}
		schematic.CUE = &v1beta1.CUE{
			Template: template,
		}
		return schematic, nil
	}

	// Check for HELM schematic
	helmData, found, err := unstructured.NestedMap(schematicData, "helm")
	if err == nil && found {
		// Handle HELM schematic
		return schematic, nil
	}

	// Check for KUBE schematic
	kubeData, found, err := unstructured.NestedMap(schematicData, "kube")
	if err == nil && found {
		// Handle KUBE schematic
		return schematic, nil
	}

	// Check for TERRAFORM schematic
	terraformData, found, err := unstructured.NestedMap(schematicData, "terraform")
	if err == nil && found {
		// Handle TERRAFORM schematic
		return schematic, nil
	}

	return nil, fmt.Errorf("no supported schematic found")
}
