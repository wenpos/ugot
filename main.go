package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("ERROR: A path parmeter needed.")
	}
	root_path := flag.Arg(0)
	TestAndAnalyzePackageCoverage(root_path)
}
