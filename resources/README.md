# resources

This directory contains the Crossplane composite resource (XR) type definitions.
Each resource is its own independent Go module, versioned and published separately,
so consuming projects can pin and upgrade individual resources without affecting others.

## Structure

```
resources/
 ├── xdeployment/
 │   ├── types.go
 │   ├── go.mod
 │   └── go.sum
 └── xappregistration/
     ├── types.go
     ├── go.mod
     └── go.sum
```

## Modules

| Module           | Import Path                                                              |
|------------------|--------------------------------------------------------------------------|
| xdeployment      | `github.com/kerwood/crossplane-xrd-generator/resources/xdeployment`      |
| xappregistration | `github.com/kerwood/crossplane-xrd-generator/resources/xappregistration` |

## Installation

Install individual resource modules as needed:

```sh
go get github.com/kerwood/crossplane-xrd-generator/resources/xdeployment@latest
go get github.com/kerwood/crossplane-xrd-generator/resources/xappregistration@latest
```

## Defining a Resource

Each resource module contains a `types.go` file with the root XR struct and its
supporting types. Structs are annotated with kubebuilder markers for validation.

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
}

type XDeploymentStatus struct{}
```

### Required vs Optional Fields

A field is **required** when its json tag has no `omitempty`:

```go
Image string `json:"image"`          // required
Tag   string `json:"tag,omitempty"`  // optional
```

### Kubebuilder Markers

Markers are standard kubebuilder validation annotations placed in comments
directly above the field they apply to. Common examples:

```go
// +kubebuilder:validation:MinLength=1
// +kubebuilder:validation:MaxLength=63
// +kubebuilder:validation:Pattern=`^[a-z0-9-]+$`
// +kubebuilder:validation:Minimum=1
// +kubebuilder:validation:Maximum=100
// +kubebuilder:validation:Enum=TCP;UDP;SCTP
// +kubebuilder:default=latest
// +kubebuilder:validation:MinItems=1
// +kubebuilder:validation:XValidation:rule="self.min <= self.max",message="min must not exceed max"
```

See the full list in the
[kubebuilder marker documentation](https://book.kubebuilder.io/reference/markers/crd-validation).

## Adding a New Resource

1. Create a new directory under `resources/`:

```sh
mkdir resources/xdatabase
```

2. Initialise a new Go module:

```sh
cd resources/xdatabase
go mod init github.com/kerwood/crossplane-xrd-generator/resources/xdatabase
go get k8s.io/apimachinery
```

3. Create `types.go` with your struct definitions and kubebuilder markers.

4. Tag and publish the new module:

```sh
git add resources/xdatabase
git commit -m "feat: add xdatabase resource"
git tag resources/xdatabase/v1.0.0
git push origin resources/xdatabase/v1.0.0
```

5. Add it to your generator tool's `go.mod` and `deps.go`.

## Versioning

Each resource module is versioned independently using Git tags with the format:

```
resources/<name>/v<major>.<minor>.<patch>
```

For example:
```sh
git tag resources/xdeployment/v1.1.0
git tag resources/xappregistration/v1.0.0
```

This allows `xdeployment` to be on `v1.1.0` while `xappregistration` stays on `v1.0.0`.

## Reusing Structs in a Crossplane Function

Since these are plain Go structs, they can be imported directly into a
[Crossplane Composite Function](https://docs.crossplane.io/latest/guides/write-a-composition-function-in-go/)
to deserialize the observed XR into a strongly typed struct:

```go
import "github.com/kerwood/crossplane-xrd-generator/resources/xdeployment"

var xr xdeployment.XDeployment
if err := json.Unmarshal(observed.Composite.Resource.Raw, &xr); err != nil {
    // handle error
}

// Access type-safe fields
image := xr.Spec.Image
replicas := xr.Spec.Replicas
```

