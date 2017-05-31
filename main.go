package main

import (
	"flag"
	"log"
)

var xerror bool
var xpkg string

func main() {
	root_path := flag.String("path", "", "[required] package path")
	ignore_error := flag.Bool("xerror", false, "[true] ignore failed test case, [false] not ignore failed test case")
	exclude_package := flag.String("xpkg", "", "exclude package name")
	flag.Parse()
	if *root_path == "" {
		log.Fatalf("ERROR: A package path needed, use --help to cat ugot parmeters")
	}
	xerror = *ignore_error
	xpkg = *exclude_package
	TestAndAnalyzePackageCoverage(*root_path)
}
