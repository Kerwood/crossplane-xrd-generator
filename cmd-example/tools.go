//go:build tools

package main

import (
	_ "github.com/kerwood/crossplane-xrd-generator/cmd-example/resources/xexample"
	_ "github.com/kerwood/crossplane-xrd-generator/resources/xappregistration"
	_ "github.com/kerwood/crossplane-xrd-generator/resources/xdeployment"
)
