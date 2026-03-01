# Crossplane XRD Generator

Instead of handcrafting the OpenAPI schema in your CompositeResourceDefinition like a caveman,
you define your composite resource (XR) as Go structs with kubebuilder markers, and this tool generates the XRD for you automatically.

The real advantage of defining your XR as Go structs is type reuse.
If you're writing a Go Function for your composition, you can deserialize the observed XR resource directly into the same
Go structs that were used to generate the XRD.

In short: define once, generate everywhere. Your XRs become type-safe and maintainable, with zero hand-crafted OpenAPI YAML to maintain.

## How It Works

Resource types are defined as Go structs in their own module and annotated with standard
[kubebuilder markers](https://book.kubebuilder.io/reference/markers/crd-validation) for validation.
The generator uses [controller-tools](https://github.com/kubernetes-sigs/controller-tools) to parse
the struct definitions and markers directly from source, producing a fully validated OpenAPI v3 schema.

This means you get full support for:
- Default values (`+kubebuilder:default`)
- Enum constraints (`+kubebuilder:validation:Enum`)
- Min/max for numbers and strings (`+kubebuilder:validation:Minimum`, `+kubebuilder:validation:MaxLength`, ...)
- CEL validation rules (`+kubebuilder:validation:XValidation`)
- And all other kubebuilder validation markers

## Defining a Resource

Create a Go module for your resource types and define your XR structs with kubebuilder markers.

```go
package xdeployment

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +kubebuilder:object:root=true
type XDeployment struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    Spec              XDeploymentSpec   `json:"spec"`
    Status            XDeploymentStatus `json:"status,omitempty"`
}

type XDeploymentSpec struct {
    // +kubebuilder:validation:MinLength=1
    Image string `json:"image"`

    // +kubebuilder:default=latest
    Tag string `json:"tag,omitempty"`

    // +kubebuilder:default=1
    // +kubebuilder:validation:Minimum=1
    // +kubebuilder:validation:Maximum=50
    Replicas *int32 `json:"replicas,omitempty"`

    // +kubebuilder:validation:Minimum=1
    // +kubebuilder:validation:Maximum=65535
    Port *int32 `json:"port,omitempty"`

    // +kubebuilder:validation:Pattern=`^[a-z0-9][a-z0-9\-\.]+[a-z0-9]$`
    Hostname *string `json:"hostname,omitempty"`
}

type XDeploymentStatus struct{}
```

Fields without `omitempty` in the json tag are required. Fields with `omitempty` are optional.

## Using the Generator

Install the generator library:

```sh
go get github.com/kerwood/crossplane-xrd-generator/generator
```

Use `ResourceMeta` to describe your resource and call `BuildCompositeResourceDefinition`:

```go
package main

import (
    "fmt"
    "os"

    "github.com/kerwood/crossplane-xrd-generator/generator"
)

func main() {
    resource := generator.ResourceMeta{
        PackagePath: "github.com/yourorg/yourrepo/resources/xdeployment",
        TypeName:    "XDeployment",
        Group:       "example.org",
        Version:     "v1alpha1",
    }

    xrd, err := generator.BuildCompositeResourceDefinition(resource)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    out, err := generator.MarshalXRDToYAML(xrd)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    os.Stdout.Write(out)
}
```

The generator resolves the package source from the module cache automatically, no local file setup required.

## Example Output

```yaml
apiVersion: apiextensions.crossplane.io/v2
kind: CompositeResourceDefinition
metadata:
  name: xdeployments.example.org
spec:
  group: example.org
  names:
    kind: XDeployment
    plural: xdeployments
  scope: Namespaced
  versions:
  - name: v1alpha1
    referenceable: true
    served: true
    schema:
      openAPIV3Schema:
        properties:
          spec:
            properties:
              image:
                type: string
                minLength: 1
              tag:
                type: string
                default: latest
              replicas:
                type: integer
                default: 1
                minimum: 1
                maximum: 50
              port:
                type: integer
                minimum: 1
                maximum: 65535
              hostname:
                type: string
                pattern: ^[a-z0-9][a-z0-9\-\.]+[a-z0-9]$
            required:
            - image
            type: object
        type: object
```

## cmd-example

For a full working example including a CLI tool that generates XRDs for multiple resources,
check out the [cmd-example](./cmd-example/) folder.

It demonstrates how to structure a project with multiple resource modules,
version them independently, and expose them through a simple CLI:

```sh
go install github.com/kerwood/crossplane-xrd-generator/cmd-example@latest

cmd-example -resource xdeployment
cmd-example -resource xappregistration
cmd-example -resource all
```

## Crossplane Function

If you are writing a [Crossplane Composite Function in Go](https://docs.crossplane.io/latest/guides/write-a-composition-function-in-go/),
you can import your XR Go structs and deserialize the observed composite resource directly into strongly typed structs.

This allows you to reuse the same structs you used to generate your XRD and work with type-safe fields instead of unstructured maps.
