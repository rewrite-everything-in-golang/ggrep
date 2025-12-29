package tests

import (
	"bytes"
	"testing"

	"github.com/rewrite-everything-in-golang/ggrep/pkg/ggrep"
)

func TestSearchReader(t *testing.T) {
	content := "Hello World\nGolang is great\nParallel search"
	reader := bytes.NewReader([]byte(content))

	config := &ggrep.Config{
		Pattern: "Golang",
	}
	err := ggrep.CompilePattern(config)
	if err != nil {
		t.Fatalf("Failed to compile pattern: %v", err)
	}

	found, err := ggrep.SearchReader(reader, "test", config)
	if err != nil {
		t.Fatalf("SearchReader failed: %v", err)
	}

	if !found {
		t.Error("Expected to find 'Golang' in content")
	}
}

func TestSearchReaderNoMatch(t *testing.T) {
	content := "Hello World\nGolang is great\nParallel search"
	reader := bytes.NewReader([]byte(content))

	config := &ggrep.Config{
		Pattern: "Rust",
	}
	err := ggrep.CompilePattern(config)
	if err != nil {
		t.Fatalf("Failed to compile pattern: %v", err)
	}

	found, err := ggrep.SearchReader(reader, "test", config)
	if err != nil {
		t.Fatalf("SearchReader failed: %v", err)
	}

	if found {
		t.Error("Expected NOT to find 'Rust' in content")
	}
}

func TestSearchReaderIgnoreCase(t *testing.T) {
	content := "Hello World\nGolang is great\nParallel search"
	reader := bytes.NewReader([]byte(content))

	config := &ggrep.Config{
		Pattern:    "golang",
		IgnoreCase: true,
	}
	err := ggrep.CompilePattern(config)
	if err != nil {
		t.Fatalf("Failed to compile pattern: %v", err)
	}

	found, err := ggrep.SearchReader(reader, "test", config)
	if err != nil {
		t.Fatalf("SearchReader failed: %v", err)
	}

	if !found {
		t.Error("Expected to find 'golang' with ignore case")
	}
}

func TestSearchReaderInvertMatch(t *testing.T) {
	content := "Relevant\nIrrelevant"
	reader := bytes.NewReader([]byte(content))

	config := &ggrep.Config{
		Pattern:     "Relevant",
		InvertMatch: true,
		LineRegexp:  true,
	}
	err := ggrep.CompilePattern(config)
	if err != nil {
		t.Fatalf("Failed to compile pattern: %v", err)
	}

	found, err := ggrep.SearchReader(reader, "test", config)
	if err != nil {
		t.Fatalf("SearchReader failed: %v", err)
	}

	if !found {
		t.Error("Expected to find non-matching lines")
	}
}

func TestSearchReaderEmpty(t *testing.T) {
	reader := bytes.NewReader([]byte(""))

	config := &ggrep.Config{
		Pattern: "any",
	}
	err := ggrep.CompilePattern(config)
	if err != nil {
		t.Fatalf("Failed to compile pattern: %v", err)
	}

	found, err := ggrep.SearchReader(reader, "test", config)
	if err != nil {
		t.Fatalf("SearchReader failed: %v", err)
	}

	if found {
		t.Error("Expected NOT to find anything in empty reader")
	}
}
