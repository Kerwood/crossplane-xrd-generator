package generator

import (
	"reflect"
	"strings"

	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	typeMetaType   = reflect.TypeOf(metav1.TypeMeta{})
	objectMetaType = reflect.TypeOf(metav1.ObjectMeta{})
)

// GoTypeToOpenAPISchema recursively converts a Go reflect.Type into a Kubernetes
// OpenAPI v3 JSONSchemaProps object.
//
// It handles:
//
//   - Primitive types: string, integer (all int kinds), boolean
//   - Slices: generates "array" with items schema
//   - Structs: generates "object" with nested properties, skipping unexported fields
//     and top-level Kubernetes metadata (TypeMeta, ObjectMeta)
//   - Recursively processes nested structs and slices
//
// The function respects `json` struct tags for property names and `required:"true"`
// tags to mark required fields. Fields with `json:"-"` are ignored. Any unknown
// type defaults to a "string" schema.
//
// This is used to automatically generate the OpenAPI schema for a resource Spec
// when building Crossplane CompositeResourceDefinitions (XRDs) or CRDs.
func GoTypeToOpenAPISchema(t reflect.Type) extv1.JSONSchemaProps {
	switch t.Kind() {

	case reflect.String:
		return extv1.JSONSchemaProps{Type: "string"}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return extv1.JSONSchemaProps{Type: "integer"}

	case reflect.Bool:
		return extv1.JSONSchemaProps{Type: "boolean"}

	case reflect.Slice:
		elemSchema := GoTypeToOpenAPISchema(t.Elem())
		return extv1.JSONSchemaProps{
			Type: "array",
			Items: &extv1.JSONSchemaPropsOrArray{
				Schema: &elemSchema,
			},
		}

	case reflect.Struct:
		props := make(map[string]extv1.JSONSchemaProps)
		required := []string{}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			// skip unexported fields
			if field.PkgPath != "" {
				continue
			}

			// skip top-level Kubernetes metadata
			if field.Type == typeMetaType || field.Type == objectMetaType {
				continue
			}

			// get json tag name
			name := field.Tag.Get("json")
			if name == "-" {
				continue
			}
			name = strings.Split(name, ",")[0]
			if name == "" {
				name = strings.ToLower(field.Name)
			}

			// recursively generate schema for this field
			prop := GoTypeToOpenAPISchema(field.Type)

			// mark required if tag says so
			if field.Tag.Get("required") == "true" {
				required = append(required, name)
			}

			props[name] = prop
		}

		return extv1.JSONSchemaProps{
			Type:       "object",
			Properties: props,
			Required:   required,
		}

	default:
		// fallback to string for unknown types
		return extv1.JSONSchemaProps{Type: "string"}
	}
}
