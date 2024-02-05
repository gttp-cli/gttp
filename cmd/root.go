package cmd

import (
	"fmt"
	"github.com/gttp-cli/gttp/pkg/model"
	"github.com/gttp-cli/gttp/pkg/parser"
	"github.com/gttp-cli/gttp/pkg/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	clip "golang.design/x/clipboard"
	"os"
)

func init() {
	rootCmd.Flags().StringP("url", "u", "", "Fetch template from URL")
	rootCmd.Flags().StringP("file", "f", "", "Fetch template from file")
	rootCmd.Flags().StringP("output", "o", "", "Output file")
	rootCmd.Flags().BoolP("clipboard", "c", false, "Copy output to clipboard")
	rootCmd.Flags().BoolP("silent", "s", false, "Silent mode")
	rootCmd.Flags().BoolP("debug", "d", false, "Print debug information")
}

var rootCmd = &cobra.Command{
	Use:   "gttp",
	Short: `Go Text Template Parser`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := cmd.Flags().GetString("url")
		file, _ := cmd.Flags().GetString("file")
		output, _ := cmd.Flags().GetString("output")
		silent, _ := cmd.Flags().GetBool("silent")
		clipboard, _ := cmd.Flags().GetBool("clipboard")
		debug, _ := cmd.Flags().GetBool("debug")

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

		tmpl, err = parser.ParseTemplate(tmpl)
		if err != nil {
			return err
		}

		result, err := parser.RenderTemplate(tmpl)
		if err != nil {
			return err
		}

		if output != "" {
			err := os.WriteFile(output, []byte(result), 0644)
			if err != nil {
				return err
			}
		}

		if clipboard {
			clip.Write(clip.FmtText, []byte(result))
		}

		if !silent {
			fmt.Println()      // padding
			fmt.Println("---") // padding
			fmt.Println()      // padding
			fmt.Println(result)
		}

		if debug {
			pterm.DefaultSection.Println("Debug Information")
			pterm.DefaultSection.WithLevel(2).Println("Template")
			pterm.Printfln("%#v", tmpl)
		}

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
