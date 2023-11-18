//go:build ignore
// +build ignore

package main

import (
	"github.com/happy-sdk/nfcsdk/internal/generator"
)

func main() {
	gen, err := generator.NewGenerator()
	if err != nil {
		generator.Fatal(err.Error())
	}
	defer gen.Finalize()

	// Run generators the exit with exit code 1
	// if something during execturin fails
	gen.Generate()
}
