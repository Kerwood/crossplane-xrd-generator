package generator

import (
	"fmt"
	"os"

	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-tools/pkg/crd"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

func ExtractOpenAPISchema(packagePath, typeName string) (*extv1.JSONSchemaProps, error) {
	roots, err := loader.LoadRoots(packagePath)
	if err != nil {
		return nil, fmt.Errorf("loading package %q: %w", packagePath, err)
	}

	if len(roots) == 0 {
		return nil, fmt.Errorf("no packages found for path %q", packagePath)
	}

	reg := &markers.Registry{}
	gen := crd.Generator{}
	if err := gen.RegisterMarkers(reg); err != nil {
		return nil, fmt.Errorf("registering markers: %w", err)
	}

	parser := &crd.Parser{
		Collector: &markers.Collector{Registry: reg},
		Checker:   &loader.TypeChecker{},
	}
	crd.AddKnownTypes(parser)

	for _, root := range roots {
		parser.NeedPackage(root)
	}

	typeIdent := crd.TypeIdent{Package: roots[0], Name: typeName}

	// Generate the raw schema first
	parser.NeedSchemaFor(typeIdent)

	// Check it exists before attempting to flatten
	if _, ok := parser.Schemata[typeIdent]; !ok {
		return nil, fmt.Errorf("type %q not found in package %q — check TypeName and PackagePath are correct", typeName, packagePath)
	}

	// Now flatten — resolves $refs and removes allOf
	parser.NeedFlattenedSchemaFor(typeIdent)

	schema, ok := parser.FlattenedSchemata[typeIdent]
	if !ok {
		return nil, fmt.Errorf("flattened schema not found for type %q in package %q", typeName, packagePath)
	}

	// Print any package errors so they're visible rather than silently causing nils
	for _, root := range roots {
		for _, err := range root.Errors {
			fmt.Fprintf(os.Stderr, "package error: %v\n", err)
		}
	}

	return &schema, nil
}
