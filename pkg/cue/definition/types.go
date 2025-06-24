package definition

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// WorkloadSchematic defines the schematic for a workload
type WorkloadSchematic struct {
	// CUE defines the CUE template
	CUE *CUE `json:"cue,omitempty"`
	// HELM defines the Helm template
	HELM *HELM `json:"helm,omitempty"`
	// Kube defines the Kubernetes template
	Kube *Kube `json:"kube,omitempty"`
}

// CUE defines the CUE template
type CUE struct {
	// Template is the CUE template string
	Template string `json:"template,omitempty"`
}

// HELM defines the Helm template
type HELM struct {
	// Release is the name of the Helm release
	Release string `json:"release,omitempty"`
	// Repository is the repository URL for the Helm chart
	Repository string `json:"repository,omitempty"`
	// Chart is the name of the Helm chart
	Chart string `json:"chart,omitempty"`
	// Version is the version of the Helm chart
	Version string `json:"version,omitempty"`
}

// Kube defines the Kubernetes template
type Kube struct {
	// Template is the Kubernetes YAML template
	Template string `json:"template,omitempty"`
	// Parameters is the list of parameters for the template
	Parameters []Parameter `json:"parameters,omitempty"`
}

// Parameter defines a parameter for a template
type Parameter struct {
	// Name is the name of the parameter
	Name string `json:"name,omitempty"`
	// Required indicates if the parameter is required
	Required bool `json:"required,omitempty"`
	// Default is the default value for the parameter
	Default interface{} `json:"default,omitempty"`
	// Description is the description of the parameter
	Description string `json:"description,omitempty"`
}

// DefinitionRevision is the base interface for all definition revisions
type DefinitionRevision interface {
	runtime.Object
	metav1.Object
	GetSpec() interface{}
}
