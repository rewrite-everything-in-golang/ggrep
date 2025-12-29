package main

import (
	"fmt"
	"os"

	"github.com/rewrite-everything-in-golang/ggrep/pkg/ggrep"
)

func main() {
	config := ggrep.ParseArgs()

	if err := ggrep.ValidateConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	if err := ggrep.CompilePattern(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	found := ggrep.SearchFilesParallel(config.Files, config)

	if found {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
