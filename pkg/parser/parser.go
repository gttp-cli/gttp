package parser

import (
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/expr-lang/expr"
	"github.com/gttp-cli/gttp/pkg/model"
	"github.com/pterm/pterm"
	"strconv"
	"strings"
	"text/template"
)

// ParseTemplate parses the template and updates its variables with filled values.
func ParseTemplate(template model.Template) (model.Template, error) {
	// Validate the template
	validationErrors := template.Validate()
	if validationErrors != nil {
		var errors []string
		for _, err := range validationErrors {
			errors = append(errors, fmt.Sprintf("- %s", err))
		}
		return template, fmt.Errorf("template validation failed:\n\n%s", strings.Join(errors, "\n"))
	}

	for i, variable := range template.Variables {
		if variable.Value != nil {
			continue // Skip variables that already have a value set.
		}

		var err error
		template.Variables[i].Value, err = processVariable(variable, template)
		if err != nil {
			return template, err
		}
	}

	return template, nil
}

func processVariable(variable model.Variable, template model.Template) (any, error) {
	if variable.Condition != "" && !evaluateCondition(variable.Condition, template) {
		return nil, nil // Condition not met, skip variable.
	}

	if strings.HasSuffix(variable.Type, "[]") {
		variable.IsArray = true
		variable.Type = strings.TrimSuffix(variable.Type, "[]")
	}

	if variable.IsArray {
		return processArrayVariable(variable, template)
	}

	return processSingleVariable(variable, template)
}

func evaluateCondition(condition string, template model.Template) bool {
	exp, err := expr.Compile(condition)
	if err != nil {
		return false
	}

	variableValues := extractVariableValues(template)
	result, err := expr.Run(exp, variableValues)
	if err != nil {
		return false
	}

	return result == true
}

func processArrayVariable(variable model.Variable, template model.Template) ([]interface{}, error) {
	var values []interface{}
	for {
		val, err := askForVariableValue(variable, template)
		if err != nil {
			return nil, err
		}

		values = append(values, val)
		if !AskToContinue() {
			break
		}
	}
	return values, nil
}

func processSingleVariable(variable model.Variable, template model.Template) (interface{}, error) {
	return askForVariableValue(variable, template)
}

func askForVariableValue(variable model.Variable, template model.Template) (any, error) {
	if structVars, ok := template.Structures[variable.Type]; ok {
		return ParseCustomType(variable, structVars)
	}
	return AskForInput(variable, "")
}

func extractVariableValues(template model.Template) map[string]interface{} {
	values := make(map[string]interface{})
	for _, variable := range template.Variables {
		values[variable.Name] = variable.Value
	}
	return values
}

func AskToContinue() bool {
	res, _ := pterm.DefaultInteractiveConfirm.Show("Add more?")
	return res
}

// AskForInput simulates asking the user for input based on the variable type and description.
func AskForInput(variable model.Variable, prefix string) (any, error) {
	var input any
	var err error

	var prompt string

	if prefix != "" {
		prompt += fmt.Sprintf("[%s] ", prefix)
	}

	if variable.Description != "" {
		prompt += variable.Description
	} else {
		prompt += variable.Name
	}

	switch variable.Type {
	case "text":
		input, err = pterm.DefaultInteractiveTextInput.WithMultiLine(variable.Multiline).Show(prompt)
		if input == "" {
			input = nil
		}
	case "number":
		var number float64
		var answer string
		answer, err = pterm.DefaultInteractiveTextInput.Show(prompt)
		if answer != "" {
			number, err = strconv.ParseFloat(answer, 64)
			input = number
		}
	case "section":
		pterm.DefaultSection.Println(variable.Name)
	case "boolean":
		input, err = pterm.DefaultInteractiveConfirm.Show(prompt)
	case "select":
		var options []string
		for _, option := range variable.Options {
			options = append(options, option.Name)
		}

		defaultOption := fmt.Sprint(variable.Default)
		input, err = pterm.DefaultInteractiveSelect.WithOptions(options).WithDefaultOption(defaultOption).Show(prompt)

		// look if option has a value
		for _, option := range variable.Options {
			if option.Name == input {
				if option.Value != nil {
					input = option.Value
				} else {
					input = option.Name
				}
			}
		}

	case "multiselect":
		var options []string
		for _, option := range variable.Options {
			options = append(options, option.Name)
		}

		var defaultOptions []string
		if variable.Default != nil {
			// If default is a string slice, add them to the default options
			// If default is a string, split it by ";" and add them to the default options
			switch variable.Default.(type) {
			case []any:
				for _, option := range variable.Default.([]any) {
					defaultOptions = append(defaultOptions, fmt.Sprint(option))
				}
			case string:
				defaultOptions = strings.Split(variable.Default.(string), ";")
			default:
				return nil, fmt.Errorf("invalid default type for multiselect: %T", variable.Default)
			}
		}
		if len(defaultOptions) == 0 {
			defaultOptions = nil
		}
		input, err = pterm.DefaultInteractiveMultiselect.WithOptions(options).WithDefaultOptions(defaultOptions).Show(prompt)
	default:
		return nil, fmt.Errorf("invalid variable type: %s", variable.Type)
	}

	if input == nil && variable.Default != nil {
		input = variable.Default
	}

	return input, err
}

// ParseCustomType handles parsing of custom types by asking for input for each field of the custom type.
func ParseCustomType(variable model.Variable, customType []model.Variable) (interface{}, error) {
	customValue := make(map[string]interface{})
	var err error
	for _, field := range customType {
		customValue[field.Name], err = AskForInput(field, variable.Name)
		if err != nil {
			return nil, err
		}
	}
	return customValue, nil
}

func ParseGoTextTemplate(templateContent string, variables map[string]any) (string, error) {
	tmpl, err := template.New("template").Funcs(sprig.FuncMap()).Parse(templateContent)
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

func RenderTemplate(template model.Template) (string, error) {
	variableValues := extractVariableValues(template)
	return ParseGoTextTemplate(template.Template, variableValues)
}
