package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kerwood/crossplane-xrd-generator/generator"
)

var xResources = map[string]generator.ResourceMeta{
	"xexample": {
		PackagePath: "github.com/kerwood/crossplane-xrd-generator/cmd-example/resources/xexample",
		TypeName:    "XExample",
		Group:       "example.org",
		Version:     "v1alpha1",
	},
	"xdeployment": {
		PackagePath: "github.com/kerwood/crossplane-xrd-generator/resources/xdeployment",
		TypeName:    "XDeployment",
		Group:       "example.org",
		Version:     "v1alpha1",
	},
	"xappregistration": {
		PackagePath: "github.com/kerwood/crossplane-xrd-generator/resources/xappregistration",
		TypeName:    "XAppRegistration",
		Group:       "example.org",
		Version:     "v1alpha1",
	},
}

func main() {
	resource := flag.String("resource", "", "XRD resource to print")
	flag.Parse()

	if *resource == "" {
		fmt.Println("Provide a resource name to print the XRD for that resource, or use 'all' to print all XRDs.")
		fmt.Println()
		fmt.Println("  Resource list:")
		for k := range xResources {
			fmt.Printf("   - %s\n", k)
		}
		fmt.Println()
		fmt.Println("  Eg. -resource xdeployment")
		os.Exit(1)
	}

	if *resource == "all" {
		for _, v := range xResources {
			printXRD(v)
			fmt.Println("---")
		}
	} else {
		meta, ok := xResources[*resource]
		if !ok {
			fmt.Printf("Error: resource '%s' not found\n", *resource)
			fmt.Println()
			fmt.Println("  Resources available:")
			for k := range xResources {
				fmt.Printf("   - %s\n", k)
			}
			os.Exit(1)
		}
		printXRD(meta)
	}
}

func printXRD(resource generator.ResourceMeta) {
	xrd, err := generator.BuildCompositeResourceDefinition(resource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building XRD: %v\n", err)
		os.Exit(1)
	}

	out, err := generator.MarshalXRDToYAML(xrd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling XRD: %v\n", err)
		os.Exit(1)
	}

	os.Stdout.Write(out)
}
