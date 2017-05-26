package main

import (
	"flag"
	"ugot/util"
	"fmt"
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("ERROR: A path parmeter needed.")
	}
	root_path := flag.Arg(0)
	util.TestAndAnalyzePackageCoverage(root_path)
}
