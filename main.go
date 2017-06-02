package main

import (
	"flag"
	"log"
)

var xerrors bool
var xpkgs string
var xfiles string

func main() {
	root_path := flag.String("path", "", "[required] package path")
	ignore_errors := flag.Bool("xerrors", false, "[true] ignore failed test case, [false] not ignore failed test case")
	exclude_packages := flag.String("xpkgs", "", "exclude packages")
	exclude_files := flag.String("xfiles", "", "exclude files")
	flag.Parse()
	if *root_path == "" {
		log.Fatalf("ERROR: A package path needed, use --help to cat ugot parmeters")
	}
	xerrors = *ignore_errors
	xpkgs = *exclude_packages
	xfiles = *exclude_files
	TestAndAnalyzePackageCoverage(*root_path)
}
