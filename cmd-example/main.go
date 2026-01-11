package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/kerwood/crossplane-xrd-generator/generator"
	"github.com/kerwood/crossplane-xrd-generator/resources/xdeployment"
	"xrd-example/resources/xexample"
)

var xResources = map[string]generator.ResourceMeta{
	"xexample": {
		Type:  reflect.TypeOf(xexample.XExample{}),
		Group: "example.org",
	},
	"xdeployment": {
		Type:  reflect.TypeOf(xdeployment.XDeployment{}),
		Group: "example.org",
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
			printXRDs(v)
			fmt.Println("---")
		}
	} else {
		rType, ok := xResources[*resource]
		if !ok {
			fmt.Printf("Error: resource '%s' not found\n", *resource)

			fmt.Println()
			fmt.Println("  Resources available:")
			for k := range xResources {
				fmt.Printf("   - %s\n", k)
			}
			os.Exit(1)
		}
		printXRDs(rType)
	}
}

func printXRDs(resource generator.ResourceMeta) {
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
