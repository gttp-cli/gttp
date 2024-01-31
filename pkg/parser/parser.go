package parser

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/pterm/pterm"
)

type Variable struct {
	Name        string
	Type        string
	IsArray     bool
	Description string

	OptionValues map[string]string // only used when type is "select"

	ComponentVars map[string]Variable

	Value        any // only set when parsed
	DefaultValue any
}

func ParseVariables(template string) ([]Variable, error) {
	variables := make([]Variable, 0)
	lines := strings.Split(template, "\n")

	var currentVar *Variable
	var componentStack []*Variable
	var inSelectOptions, inMultilineDefault bool
	var currentOption string
	var multilineDefaultValue strings.Builder

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "---" {
			if inMultilineDefault && currentVar != nil {
				currentVar.DefaultValue = multilineDefaultValue.String()
				multilineDefaultValue.Reset()
				inMultilineDefault = false
			}

			if currentVar != nil {
				variables = append(variables, *currentVar)
				currentVar = nil
			}
			break
		}

		if trimmedLine == "" {
			continue
		}

		if strings.HasPrefix(trimmedLine, "$") {
			if inMultilineDefault && currentVar != nil {
				currentVar.DefaultValue = multilineDefaultValue.String()
				multilineDefaultValue.Reset()
				inMultilineDefault = false
			}

			if inSelectOptions && currentVar != nil {
				inSelectOptions = false
				variables = append(variables, *currentVar)
				currentVar = nil
			}

			parts := strings.SplitN(trimmedLine, ":", 2)
			if len(parts) != 2 {
				continue
			}

			varName := strings.TrimSuffix(strings.TrimSpace(parts[0][1:]), "[]")
			varTypeDesc := strings.SplitN(parts[1], "//", 2)
			varType := strings.TrimSpace(varTypeDesc[0])
			varDesc := ""
			if len(varTypeDesc) == 2 {
				varDesc = strings.TrimSpace(varTypeDesc[1])
			}

			isArray := strings.HasSuffix(strings.TrimSpace(parts[0]), "[]")
			var defaultValue any

			if defaultValParts := strings.SplitN(varType, "=", 2); len(defaultValParts) == 2 {
				varType = strings.TrimSpace(defaultValParts[0])
				if varType == "boolean" {
					var b bool
					fmt.Sscan(defaultValParts[1], &b)
					defaultValue = b
				} else if varType == "number" {
					var n int
					fmt.Sscan(defaultValParts[1], &n)
					defaultValue = n
				} else {
					defaultValue = strings.TrimSpace(defaultValParts[1])
				}
			} else if strings.HasSuffix(varType, "{") {
				varType = strings.TrimSuffix(varType, " {")
				if varType == "text" {
					inMultilineDefault = true
					currentVar = &Variable{
						Name:         varName,
						Type:         varType,
						IsArray:      isArray,
						Description:  varDesc,
						OptionValues: map[string]string{},
					}
					continue
				}
			}

			newVar := Variable{
				Name:          varName,
				Type:          varType,
				IsArray:       isArray,
				Description:   varDesc,
				DefaultValue:  defaultValue,
				OptionValues:  map[string]string{},
				ComponentVars: map[string]Variable{},
			}

			if varType == "select" {
				inSelectOptions = true
				currentVar = &newVar
			} else if varType == "component" {
				if currentVar != nil && currentVar.Type == "component" {
					// Nest this component inside the current component
					currentVar.ComponentVars[varName] = newVar
					fmt.Fprintln(os.Stderr, "Nested component")
					componentStack = append(componentStack, currentVar)
				} else {
					// This is a top-level component
					fmt.Fprintln(os.Stderr, "Top-level component")
					variables = append(variables, newVar)
				}
				currentVar = &newVar
			} else {
				if currentVar != nil && currentVar.Type == "component" {
					// Add this variable as a nested variable inside the current component
					fmt.Fprintln(os.Stderr, "Nested component2", newVar.Name)
					currentVar.ComponentVars[varName] = newVar
				} else {
					// This is a top-level variable
					fmt.Fprintln(os.Stderr, "Top-level component2")
					variables = append(variables, newVar)
				}
			}
		} else if inSelectOptions {
			if trimmedLine == "{" {
				continue
			} else if trimmedLine == "}" {
				inSelectOptions = false
				variables = append(variables, *currentVar)
				for k, v := range currentVar.OptionValues {
					currentVar.OptionValues[k] = strings.TrimSpace(v)
				}
				currentVar = nil
			} else if strings.HasPrefix(line, "        ") {
				optionValue := strings.TrimSpace(strings.TrimPrefix(line, "    "))
				if currentVar != nil {
					currentVar.OptionValues[currentOption] += optionValue + "\n"
				}
			} else {
				currentOption = trimmedLine
				if currentVar != nil {
					currentVar.OptionValues[currentOption] = ""
				}
			}
		} else if inMultilineDefault {
			if trimmedLine == "}" {
				inMultilineDefault = false
				if currentVar != nil {
					currentVar.DefaultValue = strings.TrimSpace(multilineDefaultValue.String())
					multilineDefaultValue.Reset()
				}
			} else {
				multilineDefaultValue.WriteString(strings.TrimPrefix(line, "    ") + "\n")
			}
		} else if trimmedLine == "}" && len(componentStack) > 0 {
			// End of a component definition
			currentVar = componentStack[len(componentStack)-1] // Pop from stack
			componentStack = componentStack[:len(componentStack)-1]
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
	//parsed, err := f.FormatSafe(templateParts[1], VariablesToMap(variables))
	//if err != nil {
	//	return "", err
	//}

	// Parse Go text template syntax
	parsed, err := ParseGoTextTemplate(templateParts[1], VariablesToMap(variables))
	if err != nil {
		return "", err
	}

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
			var options []string
			for option := range variable.OptionValues {
				options = append(options, option)
			}
			value, err = pterm.DefaultInteractiveSelect.WithOptions(options).Show(question)
		case "multiselect":
			var options []string
			for option := range variable.OptionValues {
				options = append(options, option)
			}
			value, err = pterm.DefaultInteractiveMultiselect.WithOptions(options).Show(question)
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
