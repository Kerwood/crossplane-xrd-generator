package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"golang.org/x/tools/go/packages"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-tools/pkg/crd"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

// findModuleDir locates the source directory for a package by reading the
// build info embedded in the binary and finding the module in the module cache.
func findModuleDir(packagePath string) (string, error) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "", fmt.Errorf("build info not available")
	}

	gomodcache := os.Getenv("GOMODCACHE")
	if gomodcache == "" {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("could not determine home directory: %w", err)
			}
			gopath = filepath.Join(home, "go")
		}
		gomodcache = filepath.Join(gopath, "pkg", "mod")
	}

	// Include the main module itself in the search
	allMods := append(bi.Deps, &debug.Module{
		Path:    bi.Main.Path,
		Version: bi.Main.Version,
	})

	for _, mod := range allMods {
		if mod == nil || !strings.HasPrefix(packagePath, mod.Path) {
			continue
		}

		// Handle local replace directives
		if mod.Replace != nil {
			subPath := strings.TrimPrefix(packagePath, mod.Path)
			return filepath.Join(mod.Replace.Path, subPath), nil
		}

		version := mod.Version
		if version == "" || version == "(devel)" {
			continue
		}

		// Module cache path format: $GOMODCACHE/module@version
		moduleDir := filepath.Join(gomodcache, mod.Path+"@"+version)
		return moduleDir, nil
	}

	return "", fmt.Errorf("could not find module directory for package %q — is it a dependency?", packagePath)
}

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
