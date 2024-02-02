package cmd

import (
	"fmt"
	"github.com/gttp-cli/gttp/pkg/model"
	"github.com/gttp-cli/gttp/pkg/parser"
	"github.com/gttp-cli/gttp/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	// URL flag
	rootCmd.Flags().StringP("url", "u", "", "Fetch template from URL")

	// File flag
	rootCmd.Flags().StringP("file", "f", "", "Fetch template from file")
}

var rootCmd = &cobra.Command{
	Use:   "gttp",
	Short: `Go Text Template Parser`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := cmd.Flags().GetString("url")
		file, _ := cmd.Flags().GetString("file")

		// Do not allow both URL and file flags to be set
		if url != "" && file != "" {
			return fmt.Errorf("cannot use both URL and file flags")
		}

		// Do not allow both URL and file flags to be empty
		if url == "" && file == "" {
			return fmt.Errorf("must use either URL or file flag")
		}

		var template string
		var err error

		if url != "" {
			template, err = utils.ReadURL(url)
		} else if file != "" {
			template, err = utils.ReadFile(file)
		}

		if err != nil {
			return err
		}

		tmpl, err := model.FromYAML(template)
		if err != nil {
			return err
		}

		result, err := parser.ParseTemplate(tmpl)
		if err != nil {
			return err
		}

		fmt.Println(result)

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
