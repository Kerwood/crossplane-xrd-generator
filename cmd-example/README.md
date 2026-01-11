# Example

This example demonstrates how to use the `crossplane-xrd-generator` library to generate a Crossplane CompositeResourceDefinition from Go structs

## Files

| File                          | Description                                                                   |
| ----------------------------- | ----------------------------------------------------------------------------- |
| `resources/xexample/types.go` | Defines the `XDeployment` XR struct with `Spec` and `Status` fields           |
| `main.go`                     | Demonstrates how to use the generator library when creating a simple CLI tool |

## Running the Example

```sh
go run main.go -resource <all | xexample>
```

The command accepts a single argument, either `all` or the name of the resource you want to output.
The CLI tool will then write the generated XRD YAML to stdout.
