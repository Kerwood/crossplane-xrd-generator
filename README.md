# Crossplane XRD Generator
Instead of handcrafting the OpenAPI schema in your CompositeResourceDefinition like a caveman,
you define your composite resource (XR) as Go structs, and this tool generates the XRD for you automatically.

The real advantage of defining your XR as Go structs is type reuse.
If you’re writing a Go Function for your composition, you can deserialize the observed XR resource directly into the same
Go structs that were used to generate the XRD.

In short: define once, generate everywhere. Your XRs become type-safe and maintainable, with zero hand-crafted OpenAPI YAML to maintain.

## Example

Install `apimachinery v0.34.1` if you're running `Go 1.24` .
```sh
go get k8s.io/apimachinery/pkg/apis/meta/v1@v0.34.1
```

Install dependencies.

```sh
go get github.com/kerwood/crossplane-xrd-generator/generator
go get k8s.io/apimachinery/pkg/apis/meta/v1
```

Below is an example demonstrating how to define your Crossplane XR struct, create the corresponding `CompositeResourceDefinition`, and marshal it into YAML.

Customize the `Spec` and `Status` fields to match your needs.

```go
package main

import (
	"os"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/kerwood/crossplane-xrd-generator/generator"
)

// ---[ Crossplane Resource ]---------------------------

type XDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   XDeploymentSpec   `json:"spec"`
	Status XDeploymentStatus `json:"status,omitempty"`
}

type XDeploymentSpec struct {
	Image    string   `json:"image" required:"true"`
	Port     int      `json:"port,omitempty"`
	Hostname string   `json:"hostname,omitempty"`
}

type XDeploymentStatus struct {
	Replicas int    `json:"replicas,omitempty"`
}

// -----------------------------------------------------

func main() {
	resource := generator.ResourceMeta{
		Type:  reflect.TypeOf(XDeployment{}),
		Group: "example.org",
	}

	xrd, err := generator.BuildCompositeResourceDefinition(resource)
	if err != nil {
		panic(err)
	}

	out, err := generator.MarshalXRDToYAML(xrd)
	if err != nil {
		panic(err)
	}

	os.Stdout.Write(out)
}
```


After running the application, your Crossplane XRD will be generated automatically, ready to be applied to your cluster.

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
  - name: v1
    referenceable: true
    schema:
      openAPIV3Schema:
        properties:
          spec:
            properties:
              hostname:
                type: string
              image:
                type: string
              port:
                type: integer
            required:
            - image
            type: object
          status:
            properties:
              replicas:
                type: integer
            type: object
        type: object
    served: true
```
For additional examples and inspiration, check out the [cmd-example](./cmd-example/) folder.

## Crossplane Function
If you are writing a [Crossplane Composite Function in Go](https://docs.crossplane.io/latest/guides/write-a-composition-function-in-go/),
you can now import your XR Go struct and deserialize the observed composite resource directly into a strongly typed Go struct.

This allows you to reuse the same structs you used to generate your XRD and work with type-safe fields instead of unstructured maps.
