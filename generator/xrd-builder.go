package generator

import (
	"encoding/json"
	"reflect"
	"strings"

	apiextensionsv2 "github.com/crossplane/crossplane/v2/apis/apiextensions/v2"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceMeta holds metadata for a resource type, including its Go type
// and the API group used when building a CompositeResourceDefinition.
type ResourceMeta struct {
	Type  reflect.Type
	Group string
}

// BuildCompositeResourceDefinition creates a namespaced XRD from a Go resource type,
// generating kind, plural, and OpenAPI schema. Returns the XRD or an error.
func BuildCompositeResourceDefinition(resource ResourceMeta) (*apiextensionsv2.CompositeResourceDefinition, error) {
	schema := GoTypeToOpenAPISchema(resource.Type)

	rawSchema, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	kind := resource.Type.Name()
	plural := strings.ToLower(kind) + "s"

	return &apiextensionsv2.CompositeResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apiextensions.crossplane.io/v2",
			Kind:       "CompositeResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: plural + "." + resource.Group,
		},
		Spec: apiextensionsv2.CompositeResourceDefinitionSpec{
			Group: resource.Group,
			Scope: apiextensionsv2.CompositeResourceScopeNamespaced,

			Names: apiextv1.CustomResourceDefinitionNames{
				Kind:   kind,
				Plural: plural,
			},

			Versions: []apiextensionsv2.CompositeResourceDefinitionVersion{
				{
					Name:          "v1",
					Served:        true,
					Referenceable: true,
					Schema: &apiextensionsv2.CompositeResourceValidation{
						OpenAPIV3Schema: runtime.RawExtension{
							Raw: rawSchema,
						},
					},
				},
			},
		},
	}, nil
}
