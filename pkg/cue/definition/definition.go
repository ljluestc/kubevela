package definition

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Definition represents a generic OAM definition object.
type Definition struct {
	unstructured.Unstructured
}
