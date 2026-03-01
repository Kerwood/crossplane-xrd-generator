package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-tools/pkg/crd"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

// goModCache returns the Go module cache directory.
//
// It checks in order:
//  1. The GOMODCACHE environment variable
//  2. The output of "go env GOMODCACHE" — this handles non-standard Go
//     installations where the cache is not under ~/go/pkg/mod
//  3. A hardcoded fallback of ~/go/pkg/mod
func goModCache() string {
	if gomodcache := os.Getenv("GOMODCACHE"); gomodcache != "" {
		return gomodcache
	}

	out, err := exec.Command("go", "env", "GOMODCACHE").Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}

	home, _ := os.UserHomeDir()
	return filepath.Join(home, "go", "pkg", "mod")
}

// findModuleDir resolves the on-disk source directory for a given Go package path.
//
// It uses three strategies in order:
//
//  1. Build info — reads the module list embedded in the binary at compile time
//     via debug.ReadBuildInfo. For modules with a replace directive, the local
//     replacement path is resolved to an absolute path. For versioned modules,
//     the path is constructed from the module cache.
//
//  2. go.mod replace directives — for modules replaced with local paths that
//     do not appear in build info (e.g. because they were excluded by the linker),
//     the go.mod file in the current working directory is parsed directly.
//
//  3. go.mod require directives — for any remaining modules not matched by the
//     above, the version is read from go.mod and the module cache path is
//     constructed as $GOMODCACHE/<module>@<version>.
func findModuleDir(packagePath string) (string, error) {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, mod := range bi.Deps {
			if mod == nil || !strings.HasPrefix(packagePath, mod.Path) {
				continue
			}
			if mod.Replace != nil {
				subPath := strings.TrimPrefix(packagePath, mod.Path)
				absPath, err := filepath.Abs(mod.Replace.Path)
				if err != nil {
					return "", fmt.Errorf("resolving replace path: %w", err)
				}
				return filepath.Join(absPath, subPath), nil
			}
			if mod.Version != "" && mod.Version != "(devel)" {
				return filepath.Join(goModCache(), mod.Path+"@"+mod.Version), nil
			}
		}
	}

	// Fall back to parsing go.mod directly for replaced modules
	gomodBytes, err := os.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("reading go.mod: %w", err)
	}

	mf, err := modfile.Parse("go.mod", gomodBytes, nil)
	if err != nil {
		return "", fmt.Errorf("parsing go.mod: %w", err)
	}

	for _, r := range mf.Replace {
		if !strings.HasPrefix(packagePath, r.Old.Path) {
			continue
		}
		subPath := strings.TrimPrefix(packagePath, r.Old.Path)
		absPath, err := filepath.Abs(r.New.Path)
		if err != nil {
			return "", fmt.Errorf("resolving replace path: %w", err)
		}
		return filepath.Join(absPath, subPath), nil
	}

	// Fall back to module cache lookup via require directives
	for _, r := range mf.Require {
		if !strings.HasPrefix(packagePath, r.Mod.Path) {
			continue
		}
		return filepath.Join(goModCache(), r.Mod.Path+"@"+r.Mod.Version), nil
	}

	return "", fmt.Errorf("could not find module directory for package %q", packagePath)
}

// ExtractOpenAPISchema parses a Go package by import path and returns a fully
// flattened OpenAPI v3 JSONSchemaProps for the named type.
//
// It uses controller-tools to walk the package's AST, collect kubebuilder
// marker comments (e.g. +kubebuilder:validation:Minimum, +kubebuilder:default),
// and generate a schema that reflects both the struct field types and any
// validation constraints defined via markers.
//
// The returned schema has all $ref pointers resolved and allOf wrappers removed,
// making it suitable for direct embedding in a Crossplane XRD or Kubernetes CRD
// without further processing.
//
// packagePath is the full Go import path, e.g.
// "github.com/yourorg/yourrepo/resources/xdeployment".
// typeName is the root struct name to generate the schema for, e.g. "XDeployment".
func ExtractOpenAPISchema(packagePath, typeName string) (*extv1.JSONSchemaProps, error) {
	moduleDir, err := findModuleDir(packagePath)
	if err != nil {
		return nil, fmt.Errorf("finding module dir: %w", err)
	}

	cfg := &packages.Config{Dir: moduleDir}
	roots, err := loader.LoadRootsWithConfig(cfg, packagePath)
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
	parser.NeedSchemaFor(typeIdent)

	if _, ok := parser.Schemata[typeIdent]; !ok {
		for _, root := range roots {
			for _, e := range root.Errors {
				fmt.Fprintf(os.Stderr, "package error: %v\n", e)
			}
		}
		return nil, fmt.Errorf("type %q not found in package %q", typeName, packagePath)
	}

	parser.NeedFlattenedSchemaFor(typeIdent)

	schema, ok := parser.FlattenedSchemata[typeIdent]
	if !ok {
		return nil, fmt.Errorf("flattened schema not found for type %q", typeName)
	}

	return &schema, nil
}
