package main

import (
	"github.com/edgexr/go-vet/analyzers/badfuncs"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		//shadow.Analyzer,
		badfuncs.Analyzer,
	)
}
