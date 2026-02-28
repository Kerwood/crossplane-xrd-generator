package generator

import (
	"encoding/json"
	"fmt"
	"strings"

	apiextensionsv2 "github.com/crossplane/crossplane/v2/apis/apiextensions/v2"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceMeta holds metadata for generating an XRD from a Go package.
type ResourceMeta struct {
	// PackagePath is the Go import path to the package containing the type.
	// e.g. "github.com/yourorg/yourrepo/resources/xdeployment"
	PackagePath string

	// TypeName is the root struct name to generate the schema from.
	// e.g. "XDeployment"
	TypeName string

	// Group is the Crossplane API group.
	// e.g. "example.crossplane.io"
	Group string

	// Version defaults to "v1alpha1" if empty.
	Version string
}

func BuildCompositeResourceDefinition(resource ResourceMeta) (*apiextensionsv2.CompositeResourceDefinition, error) {
	version := resource.Version
	if version == "" {
		version = "v1alpha1"
	}

	// Extract schema from Go package using controller-tools
	schema, err := ExtractOpenAPISchema(resource.PackagePath, resource.TypeName)
	if err != nil {
		return nil, fmt.Errorf("extracting schema: %w", err)
	}

	// Pull the spec properties out of the top-level schema.
	// controller-tools generates a schema for the whole struct (XDeployment),
	// but XRD only wants the spec properties.
	specSchema, ok := schema.Properties["spec"]
	if !ok {
		return nil, fmt.Errorf("no 'spec' field found in schema for %q", resource.TypeName)
	}

	rawSchema, err := json.Marshal(specSchema)
	if err != nil {
		return nil, err
	}

	kind := resource.TypeName
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
					Name:          version,
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
