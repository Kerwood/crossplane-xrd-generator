# cmd-example

A CLI tool that demonstrates how to use the `generator` package to generate
Crossplane `CompositeResourceDefinition` YAML from Go struct definitions.

## Installation

```sh
go install github.com/kerwood/crossplane-xrd-generator/cmd-example@latest
```

## Usage

```sh
# Print XRD for a specific resource
cmd-example -resource xdeployment

# Print XRDs for all resources
cmd-example -resource all
```

Running without arguments prints the available resources:

```sh
cmd-example

Provide a resource name to print the XRD for that resource, or use 'all' to print all XRDs.

  Resource list:
   - xdeployment
   - xappregistration
   - xexample

  Eg. -resource xdeployment
```

## Piping to a File

```sh
cmd-example -resource xdeployment > xdeployment.yaml
cmd-example -resource all > xrds.yaml
```

## Project Structure

```
cmd-example/
 ├── main.go       ← CLI entrypoint
 ├── deps.go       ← blank imports to keep resource modules in go.mod
 ├── go.mod
 └── go.sum
```

All resources are separate modules defined under the [resources](../resources/)
directory and imported as dependencies.

## Adding a Resource

1. If the resource is in an external module, add it to `go.mod`:

```sh
go get github.com/kerwood/crossplane-xrd-generator/resources/xnewresource@latest
```

2. Add a blank import to `deps.go` to ensure the module stays in `go.mod` after `go mod tidy`:

```go
package main

import (
    _ "github.com/kerwood/crossplane-xrd-generator/resources/xdeployment"
    _ "github.com/kerwood/crossplane-xrd-generator/resources/xappregistration"
    _ "github.com/kerwood/crossplane-xrd-generator/resources/xnewresource"  // ← add this
)
```

3. Register it in the `xResources` map in `main.go`:

```go
var xResources = map[string]generator.ResourceMeta{
    "xnewresource": {
        PackagePath: "github.com/kerwood/crossplane-xrd-generator/resources/xnewresource",
        TypeName:    "XNewResource",
        Group:       "example.org",
        Version:     "v1alpha1",
    },
    // ... existing resources
}
```

## Local Development

During development, use `replace` directives in `go.mod` to point at local source
instead of published versions, so changes are picked up immediately without
committing and tagging:

```
replace (
    github.com/kerwood/crossplane-xrd-generator/generator => ../generator
    github.com/kerwood/crossplane-xrd-generator/resources/xdeployment => ../resources/xdeployment
    github.com/kerwood/crossplane-xrd-generator/resources/xappregistration => ../resources/xappregistration
)
```

Then run directly with:

```sh
go run main.go -resource xdeployment
```

Remember to remove the `replace` directives and update the versions before publishing.
