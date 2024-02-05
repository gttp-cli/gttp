package cmd

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gttp-cli/gttp/pkg/model"
	"github.com/gttp-cli/gttp/pkg/parser"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	rootCmd.AddCommand(serveCmd)

	// Add address flag
	serveCmd.Flags().StringP("address", "a", "localhost:8080", "Address to listen on")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, _ := cmd.Flags().GetString("address")
		app := fiber.New()

		app.Use(logger.New())

		app.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(map[string]string{
				"status": "ok",
				"docs":   "https://docs.gttp.dev",
			})
		})

		api := app.Group("/api")
		v1 := api.Group("/v1")

		// /parse accepts YAML and returns the parsed template as JSON
		v1.Post("/parse", func(c *fiber.Ctx) error {
			// Get template from JSON "template" key
			body := struct {
				Template string `json:"template"`
			}{}

			if err := c.BodyParser(&body); err != nil {
				return c.Status(400).JSON(map[string]string{
					"error": err.Error(),
				})
			}

			var tmpl model.Template
			var err error

			if strings.HasPrefix(body.Template, "{") {
				tmpl, err = model.FromJSON(body.Template)
			} else {
				tmpl, err = model.FromYAML(body.Template)
			}

			if err != nil {
				return c.Status(400).JSON(map[string]string{
					"error": err.Error(),
				})
			}

			// Validate template
			errs := tmpl.Validate()
			if errs != nil {
				var errors []string
				for _, err := range errs {
					errors = append(errors, err.Error())
				}
				return c.Status(400).JSON(map[string]interface{}{
					"errors": errors,
				})
			}

			rendered, err := parser.RenderTemplate(tmpl)
			if err != nil {
				return c.Status(500).JSON(map[string]string{
					"error": err.Error(),
				})
			}

			return c.JSON(map[string]string{
				"template": body.Template,
				"rendered": rendered,
			})
		})

		return app.Listen(addr)
	},
}
