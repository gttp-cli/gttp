package cmd

import (
	"fmt"
	"github.com/gttp-cli/gttp/pkg/model"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(specCmd)

	// Add output flag
	specCmd.Flags().StringP("output", "o", "", "Output file")
}

var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "Create a template that shows the full specifcation of gttp templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		template := model.Template{
			Structures: map[string][]model.Variable{
				"person": {
					{
						Name:        "Name",
						Type:        "text",
						Description: "Name of the person",
					},
					{
						Name:        "Age",
						Type:        "number",
						Description: "Age of the person",
					},
				},
			},
			Variables: []model.Variable{
				{
					Name:        "Name",
					Type:        "text",
					Description: "Name of the user",
					Default:     "World",
				},
				{
					Name:        "Age",
					Type:        "number",
					Description: "Age of the user",
				},
				{
					Name:        "Animal",
					Type:        "select",
					Description: "Favorite animal",
					Options: []model.Option{
						{Name: "Cat", Value: "cat"},
						{Name: "Dog", Value: "dog"},
					},
				},
			},
			Template: `
Hello {{ .Name }}, you're a {{ .Animal }} person!
{{ .Age }} is a great age!`,
		}

		tmpl, err := template.ToYAML()
		if err != nil {
			return err
		}

		if output, _ := cmd.Flags().GetString("output"); output != "" {
			err = os.WriteFile(output, []byte(tmpl), 0644)
			if err != nil {
				return err
			}
		} else {
			fmt.Println(tmpl)
		}

		return nil
	},
}
