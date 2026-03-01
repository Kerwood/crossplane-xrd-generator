# generator

The `generator` package is the core library of the Crossplane XRD Generator.
It takes a Go package path and type name, parses the struct definitions and kubebuilder markers
using [controller-tools](https://github.com/kubernetes-sigs/controller-tools), and produces a
fully populated `CompositeResourceDefinition` object ready to be marshalled to YAML.

## Packages

| File | Responsibility |
|---|---|
| `schema-extractor.go` | Locates the package source on disk and extracts a flattened OpenAPI v3 schema using controller-tools |
| `xrd-builder.go` | Wraps the schema in a Crossplane `CompositeResourceDefinition` struct |
| `emitter.go` | Marshals the XRD to apply-ready YAML |

## Installation

```sh
go get github.com/kerwood/crossplane-xrd-generator/generator
```

## Usage

```go
import "github.com/kerwood/crossplane-xrd-generator/generator"

resource := generator.ResourceMeta{
    PackagePath: "github.com/yourorg/yourrepo/resources/xdeployment",
    TypeName:    "XDeployment",
    Group:       "example.org",
    Version:     "v1alpha1",
}

xrd, err := generator.BuildCompositeResourceDefinition(resource)
if err != nil {
    // handle error
}

out, err := generator.MarshalXRDToYAML(xrd)
if err != nil {
    // handle error
}

os.Stdout.Write(out)
```

## ResourceMeta

`ResourceMeta` is the input struct passed to `BuildCompositeResourceDefinition`.

| Field         | Required | Description                                                                                                      |
|---------------|----------|------------------------------------------------------------------------------------------------------------------|
| `PackagePath` | yes      | Full Go import path to the package containing the type, e.g. `github.com/yourorg/yourrepo/resources/xdeployment` |
| `TypeName`    | yes      | Name of the root struct to generate the schema from, e.g. `XDeployment`                                          |
| `Group`       | yes      | Crossplane API group, e.g. `example.org`                                                                         |
| `Version`     | no       | API version. Defaults to `v1alpha1` if empty                                                                     |

## How Schema Extraction Works

`ExtractOpenAPISchema` resolves the package source directory using the following
strategies in order:

1. **Build info** - reads the module list embedded in the binary at compile time via
   `debug.ReadBuildInfo`. For versioned modules the source is located in the module
   cache. For modules with a `replace` directive, the local replacement path is used.

2. **go.mod replace directives** - if a module is not found in build info (e.g. because
   it was replaced with a local path and excluded by the linker), the `go.mod` in the
   current working directory is parsed directly to find the replacement path.

3. **go.mod require directives** - for any remaining modules, the version is read from
   `go.mod` and the path is constructed as `$GOMODCACHE/<module>@<version>`.

The module cache location is resolved via `go env GOMODCACHE`, falling back to
`~/go/pkg/mod` if the command is unavailable.

Once the source directory is located, controller-tools loads the package AST, collects
all kubebuilder marker comments, and generates a fully flattened OpenAPI v3 schema with
all `$ref` pointers resolved and `allOf` wrappers removed.

## Supported Kubebuilder Markers

Any marker supported by controller-tools works out of the box. Common examples:

```go
// +kubebuilder:validation:MinLength=1
// +kubebuilder:validation:MaxLength=63
// +kubebuilder:validation:Pattern=`^[a-z0-9-]+$`
// +kubebuilder:validation:Minimum=1
// +kubebuilder:validation:Maximum=100
// +kubebuilder:validation:Enum=TCP;UDP;SCTP
// +kubebuilder:default=latest
// +kubebuilder:validation:MinItems=1
// +kubebuilder:validation:MaxItems=10
// +kubebuilder:validation:XValidation:rule="self.min <= self.max",message="min must not exceed max"
```

See the full list in the
[kubebuilder marker documentation](https://book.kubebuilder.io/reference/markers/crd-validation).

## Required vs Optional Fields

A field is **required** when its json tag has no `omitempty`:

```go
Image string `json:"image"`           // required
Tag   string `json:"tag,omitempty"`   // optional
```

No additional marker is needed to mark a field as required.

