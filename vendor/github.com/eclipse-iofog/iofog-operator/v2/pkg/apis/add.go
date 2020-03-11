package apis

import (
	"github.com/eclipse-iofog/iofog-operator/v2/pkg/apis/iofog/v1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1.SchemeBuilder.AddToScheme)
}
