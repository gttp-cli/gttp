package main

import (
	"encoding/json"
	"github.com/MarvinJWendt/testza"
	"github.com/gttp-cli/gttp/pkg/parser"
	"github.com/gttp-cli/gttp/pkg/utils"
	"os"
	"path/filepath"
	"testing"
)

func TestExamples(t *testing.T) {
	const exampleDir = "_examples"

	// Recursively walk through all *.gttp files in the example directory
	err := filepath.Walk(exampleDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip non-gttp files
		if filepath.Ext(path) != ".gttp" {
			return nil
		}

		t.Run(path, func(t *testing.T) {
			content, err := utils.ReadFile(path)
			if err != nil {
				t.Error(err)
			}

			t.Run("parse-variables", func(t *testing.T) {
				variables, err := parser.ParseVariables(content)
				if err != nil {
					t.Error(err)
				}

				json, err := json.MarshalIndent(variables, "", "  ")
				if err != nil {
					t.Error(err)
				}

				t.Run("compare", func(t *testing.T) {
					// Compare the output with the expected output
					expectedPath := filepath.Dir(path) + string(os.PathSeparator) + "parsed.json"

					// Check if the expected output file exists
					if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
						// Write JSON to file
						err = os.WriteFile(expectedPath, json, 0644)
						if err != nil {
							t.Error(err)
						}
					}

					expected, err := utils.ReadFile(expectedPath)
					if err != nil {
						t.Error(err)
					}

					testza.AssertEqual(t, expected, string(json))
				})
			})
		})

		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
