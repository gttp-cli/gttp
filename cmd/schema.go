package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gttp-cli/gttp/pkg/model"
	"github.com/invopop/jsonschema"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(schemaCmd)

	// Add output flag
	schemaCmd.Flags().StringP("output", "o", "", "Output file")
}

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Create JSON schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		outputPath, _ := cmd.Flags().GetString("output")

		r := new(jsonschema.Reflector)
		if err := r.AddGoComments("github.com/gttp-cli/gttp", "./"); err != nil {
			// deal with error
		}
		schema := r.Reflect(&model.Template{})

		j, err := schema.MarshalJSON()
		if err != nil {
			return fmt.Errorf("failed to generate JSON schema: %w", err)
		}

		// Format JSON
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, j, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format JSON schema: %w", err)
		}

		j = prettyJSON.Bytes()

		if outputPath != "" {
			err = os.WriteFile(outputPath, j, 0644)
			if err != nil {
				return fmt.Errorf("failed to write JSON schema to file: %w", err)
			}
		} else {
			fmt.Println(string(j))
		}

		return nil
	},
}
