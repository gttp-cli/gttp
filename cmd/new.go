package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCmd)

	// Add output flag
	newCmd.Flags().StringP("output", "o", "", "Output file")
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new template",
	RunE: func(cmd *cobra.Command, args []string) error {
		//		template := model.Template{
		//			Variables: model.Variables{
		//				Variables: []model.Variable{
		//					{
		//						Name:        "Name",
		//						Type:        "text",
		//						Description: "Name of the user",
		//						Default:     "World",
		//					},
		//					{
		//						Name:        "Age",
		//						Type:        "number",
		//						Description: "Age of the user",
		//					},
		//					{
		//						Name:        "Animal",
		//						Type:        "select",
		//						Description: "Favorite animal",
		//						Options: []model.Option{
		//							{Name: "Cat", Value: "cat"},
		//							{Name: "Dog", Value: "dog"},
		//						},
		//					},
		//				},
		//			},
		//			Template: `
		//Hello {{ .Name }}, you're a {{ .Animal }} person!
		//{{ .Age }} is a great age!`,
		//		}
		//
		//		tmpl, err := template.ToYAML()
		//		if err != nil {
		//			return err
		//		}
		//
		//		if output, _ := cmd.Flags().GetString("output"); output != "" {
		//			err = os.WriteFile(output, []byte(tmpl), 0644)
		//			if err != nil {
		//				return err
		//			}
		//		} else {
		//			fmt.Println(tmpl)
		//		}

		return nil
	},
}
