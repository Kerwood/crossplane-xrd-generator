package main

import (
	// Import resource packages to ensure they appear in build info
	// so the schema extractor can locate them in the module cache.
	_ "github.com/kerwood/crossplane-xrd-generator/resources/xappregistration"
	_ "github.com/kerwood/crossplane-xrd-generator/resources/xdeployment"
	_ "github.com/kerwood/crossplane-xrd-generator/resources/xexample"
)
