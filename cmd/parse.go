package cmd

import (
	"errors"
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/gttp-cli/gttp/pkg/model"
	"github.com/gttp-cli/gttp/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(parseCmd)

	parseCmd.Flags().StringP("url", "u", "", "Fetch template from URL")
	parseCmd.Flags().StringP("file", "f", "", "Fetch template from file")
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: `Parse template and print AST`,
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
			return errors.New(yaml.FormatError(err, true, true))
		}

		j, err := tmpl.ToJSON()
		if err != nil {
			return err
		}
		fmt.Println(j)

		return nil
	},
}
