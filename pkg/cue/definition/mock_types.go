package definition

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MockRevisionSpec represents a mock revision spec for tests
type MockRevisionSpec struct {
	// Template is the template string
	Template string `json:"template,omitempty"`
	// DefinitionType is the type of definition
	DefinitionType string `json:"definitionType,omitempty"`
}

// MockWorkflowStepDefinition represents a mock workflow step definition for tests
type MockWorkflowStepDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MockWorkflowStepSpec `json:"spec,omitempty"`
}

// MockWorkflowStepSpec represents a mock workflow step spec for tests
type MockWorkflowStepSpec struct {
	// Template is the template string
	Template string `json:"template,omitempty"`
	// Type is the type of the workflow step
	Type string `json:"type,omitempty"`
}

// MockComponentDefinitionRevision represents a mock component definition revision for tests
type MockComponentDefinitionRevision struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MockRevisionSpec `json:"spec,omitempty"`
}

// MockTraitDefinitionRevision represents a mock trait definition revision for tests
type MockTraitDefinitionRevision struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MockRevisionSpec `json:"spec,omitempty"`
}

// MockWorkflowStepDefinitionRevision represents a mock workflow step definition revision for tests
type MockWorkflowStepDefinitionRevision struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MockRevisionSpec `json:"spec,omitempty"`
}
