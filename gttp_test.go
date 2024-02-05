package main

import (
	"fmt"
	"github.com/gttp-cli/gttp/pkg/model"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTestFiles(t *testing.T) {
	// Walk over "testdata" and "_examples" and execute test for all ".yml" files.
	testPaths := []string{"testdata", "_examples"}
	for _, testPath := range testPaths {
		filepath.Walk(testPath, func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".yml" {
				t.Run(path, func(t *testing.T) {
					mustFail := strings.Contains(path, "must-fail")
					if mustFail {
						testFileMustFail(t, path)
					} else {
						testFileMustParse(t, path)
					}
				})
			}

			return nil
		})
	}
}

func testFileMustParse(t *testing.T, path string) {
	err := runTestFile(t, path)
	if err != nil {
		t.Fatal(err)
	}
}

func testFileMustFail(t *testing.T, path string) {
	err := runTestFile(t, path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func runTestFile(t *testing.T, path string) error {
	// Read file
	file, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
		return nil
	}

	template, err := model.FromYAML(string(file))
	if err != nil {
		t.Fatal(err)
		return nil
	}

	// Validate template
	errs := template.Validate()
	if len(errs) > 0 {
		var errors []string
		for _, err := range errs {
			errors = append(errors, err.Error())
		}
		return fmt.Errorf("template validation failed:\n\n%s", strings.Join(errors, "\n"))
	}

	return nil
}
