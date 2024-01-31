package parser

import (
	"atomicgo.dev/f"
	"fmt"
	"github.com/pterm/pterm"
	"strings"
	"text/template"
)

type Variable struct {
	Name        string
	Description string
	Type        string
	Options     []string // only used when type is "select"
	Value       any      // only set when parsed
}

func ParseVariables(template string) ([]Variable, error) {
	variables := make([]Variable, 0)
	lines := strings.Split(template, "\n")

	var currentVar *Variable
	for _, line := range lines {
		// Stop parsing if we reach the delimiter
		if strings.TrimSpace(line) == "---" {
			break
		}

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Check for variable definition
		if strings.HasPrefix(line, "$") {
			// Extract variable name, type, and description
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue // Invalid format, skip this line
			}

			varDescParts := strings.SplitN(parts[1], "//", 2)
			varType := strings.TrimSpace(varDescParts[0])
			varDesc := ""
			if len(varDescParts) == 2 {
				varDesc = strings.TrimSpace(varDescParts[1])
			}

			currentVar = &Variable{
				Name:        strings.TrimSpace(parts[0][1:]), // remove '$' from start
				Type:        varType,
				Description: varDesc,
				Options:     []string{},
			}
			variables = append(variables, *currentVar)
		} else if currentVar != nil && (currentVar.Type == "select" || currentVar.Type == "multiselect") {
			// Parse select or multiselect options
			option := strings.TrimSpace(line)
			if strings.HasPrefix(option, "-") {
				option = strings.TrimSpace(option[1:])
				currentVar.Options = append(currentVar.Options, option)
				// Update the last variable in the slice with new options
				variables[len(variables)-1] = *currentVar
			}
		}
	}

	return variables, nil
}

func ParseTemplate(template string) (string, error) {
	template = strings.ReplaceAll(template, "\r\n", "\n")
	templateParts := strings.SplitN(template, "\n---\n", 2)

	variables, err := ParseVariables(templateParts[0])
	if err != nil {
		return "", fmt.Errorf("failed to parse variables: %w", err)
	}

	err = AskForVariables(variables)
	if err != nil {
		return "", err
	}

	// Parse expr-lang syntax
	parsed, err := f.FormatSafe(templateParts[1], VariablesToMap(variables))
	if err != nil {
		return "", err
	}

	// Parse Go text template syntax
	parsed, err = ParseGoTextTemplate(parsed, VariablesToMap(variables))

	return parsed, nil
}

func VariablesToMap(variables []Variable) map[string]any {
	variableMap := make(map[string]any)

	for _, variable := range variables {
		variableMap[variable.Name] = variable.Value
	}

	return variableMap
}

func AskForVariables(variables []Variable) error {
	for i, variable := range variables {
		// Skip variables that already have a value
		if variable.Value != nil {
			continue
		}

		var question string

		if variable.Description != "" {
			question = variable.Description
		} else {
			question = variable.Name
		}

		question += fmt.Sprintf(pterm.Gray(" (%s)"), variable.Type)

		// Ask for variable value
		var value any
		var err error
		switch variable.Type {
		case "text", "string":
			value, err = pterm.DefaultInteractiveTextInput.Show(question)
		case "number", "int", "integer":
			value, err = pterm.DefaultInteractiveTextInput.Show(question)
		case "bool", "boolean":
			value, err = pterm.DefaultInteractiveConfirm.Show(question)
		case "select":
			value, err = pterm.DefaultInteractiveSelect.WithOptions(variable.Options).Show(question)
		case "multiselect":
			value, err = pterm.DefaultInteractiveMultiselect.WithOptions(variable.Options).Show(question)
		}
		if err != nil {
			return err
		}

		// Update variable value
		variables[i].Value = value
	}

	return nil
}

func ParseGoTextTemplate(templateContent string, variables map[string]any) (string, error) {
	// Parse Go text template syntax
	tmpl, err := template.New("template").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse go template: %w", err)
	}

	var parsed strings.Builder
	err = tmpl.Execute(&parsed, variables)
	if err != nil {
		return "", fmt.Errorf("failed to execute go template: %w", err)
	}

	return parsed.String(), nil
}
