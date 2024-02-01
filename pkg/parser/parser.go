package parser

import (
	"atomicgo.dev/f"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/pterm/pterm"
)

type Variable struct {
	Name        string
	Type        string
	IsArray     bool
	Description string

	OptionValues map[string]string // only used when type is "select" or "multiselect"
	OptionOrder  []string          // only used when type is "select" or "multiselect"

	ComponentVars  map[string]Variable
	ComponentOrder *[]string

	Value        any // only set when parsed
	DefaultValue any
}

func ParseVariables(template string) ([]Variable, error) {
	variables := make([]Variable, 0)
	lines := strings.Split(template, "\n")

	var currentVar *Variable
	var componentStack []*Variable
	var inSelectOptions, inMultilineDefault, inComponent bool
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

			if inComponent {
				inComponent = false
				if len(componentStack) > 0 {
					componentStack = componentStack[:len(componentStack)-1]
					if len(componentStack) == 0 && currentVar != nil {
						variables = append(variables, *currentVar)
						currentVar = nil
					}
				}
			}
			break
		}

		if trimmedLine == "" {
			continue
		}

		if strings.HasPrefix(trimmedLine, "# ") {
			newVar := Variable{
				Value: strings.TrimSpace(trimmedLine[2:]),
				Type:  "section",
			}
			variables = append(variables, newVar)
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
						Name:           varName,
						Type:           varType,
						IsArray:        isArray,
						Description:    varDesc,
						OptionValues:   map[string]string{},
						OptionOrder:    []string{},
						ComponentOrder: &[]string{},
					}
					continue
				}
			}

			newVar := Variable{
				Name:           varName,
				Type:           varType,
				IsArray:        isArray,
				Description:    varDesc,
				DefaultValue:   defaultValue,
				OptionValues:   map[string]string{},
				OptionOrder:    []string{},
				ComponentVars:  map[string]Variable{},
				ComponentOrder: &[]string{},
			}

			if varType == "select" || varType == "multiselect" {
				inSelectOptions = true
				currentVar = &newVar
			} else if varType == "component" {
				if inComponent {
					// If we are already in a component, add newVar as nested component
					componentStack[len(componentStack)-1].ComponentVars[varName] = newVar
				} else {
					// If we are not in a component, add newVar to variables
					variables = append(variables, newVar)
				}
				// Push newVar onto the stack and set inComponent to true
				componentStack = append(componentStack, &newVar)
				inComponent = true
				currentVar = &newVar
			} else {
				if inComponent {
					// If in a component, add newVar as nested variable
					componentStack[len(componentStack)-1].ComponentVars[varName] = newVar
					// Add component order
					*componentStack[len(componentStack)-1].ComponentOrder = append(*componentStack[len(componentStack)-1].ComponentOrder, varName)
				} else {
					// If not in a component, add newVar to variables
					variables = append(variables, newVar)
				}
			}
		} else if inSelectOptions {
			if trimmedLine == "{" {
				continue
			} else if trimmedLine == "}" {
				inSelectOptions = false
				if currentVar != nil {
					for k, v := range currentVar.OptionValues {
						currentVar.OptionValues[k] = strings.TrimSpace(v)
					}
					variables = append(variables, *currentVar)
					currentVar = nil
				}
			} else if strings.HasPrefix(line, "        ") {
				optionValue := strings.TrimSpace(strings.TrimPrefix(line, "    "))
				if currentVar != nil {
					currentVar.OptionValues[currentOption] += optionValue + "\n"
				}
			} else {
				currentOption = trimmedLine
				if currentVar != nil {
					currentVar.OptionValues[currentOption] = ""
					currentVar.OptionOrder = append(currentVar.OptionOrder, currentOption)
				}
			}
		} else if inMultilineDefault {
			if trimmedLine == "}" {
				inMultilineDefault = false
				if currentVar != nil {
					currentVar.DefaultValue = strings.TrimSpace(multilineDefaultValue.String())
					multilineDefaultValue.Reset()
					if inComponent {
						componentStack[len(componentStack)-1].ComponentVars[currentVar.Name] = *currentVar
					} else {
						variables = append(variables, *currentVar)
					}
				}
			} else {
				multilineDefaultValue.WriteString(strings.TrimPrefix(line, "    ") + "\n")
			}
		} else if trimmedLine == "}" {
			if inComponent {
				// Pop the current component off the stack
				componentStack = componentStack[:len(componentStack)-1]
				inComponent = len(componentStack) > 0
				if len(componentStack) > 0 {
					currentVar = componentStack[len(componentStack)-1]
				} else {
					currentVar = nil
				}
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

	err = AskForVariables(variables, "")
	if err != nil {
		return "", err
	}

	//Parse expr-lang syntax
	varMap := VariablesToMap(variables)

	parsed, err := f.FormatSafe(templateParts[1], varMap)
	if err != nil {
		return "", err
	}

	// Parse Go text template syntax
	parsed, err = ParseGoTextTemplate(parsed, varMap)
	if err != nil {
		return "", err
	}

	return parsed, nil
}

func VariablesToMap(variables []Variable) map[string]any {
	variableMap := make(map[string]any)

	for _, variable := range variables {
		if variable.Name == "" { // Skip if no name
			continue
		}

		// Check for nil value to avoid nil pointer panic
		if variable.Value == nil {
			continue
		}

		valueType := reflect.TypeOf(variable.Value)
		if valueType != nil && valueType.Kind() == reflect.Slice {
			// Convert slice of Variables into a nested map
			nestedMap := make(map[string]any)
			s := reflect.ValueOf(variable.Value)
			for i := 0; i < s.Len(); i++ {
				nestedVar := s.Index(i).Interface().(Variable)
				if nestedVar.Value != nil {
					nestedMap[nestedVar.Name] = nestedVar.Value
				}
			}
			variableMap[variable.Name] = nestedMap
		} else {
			variableMap[variable.Name] = variable.Value
		}
	}

	return variableMap
}

func AskForVariables(variables []Variable, path string) error {
	for i, variable := range variables {
		if variable.Type == "section" {
			pterm.DefaultSection.Println(variable.Value)
			continue
		}

		// Skip variables that already have a value
		if variable.Value != nil {
			continue
		}

		var question string

		if path != "" {
			question += path + "."
		}

		if variable.Description != "" {
			question += variable.Description
		} else {
			question += variable.Name
		}
		question += fmt.Sprintf(pterm.Gray(" (%s)"), variable.Type)

		// Ask for variable value
		var value any
		var err error
		switch variable.Type {
		case "text":
			value, err = pterm.DefaultInteractiveTextInput.Show(question)
		case "number":
			value, err = pterm.DefaultInteractiveTextInput.Show(question)
		case "boolean":
			value, err = pterm.DefaultInteractiveConfirm.Show(question)
		case "select":
			var options []string
			for _, option := range variable.OptionOrder {
				options = append(options, option)
			}
			value, err = pterm.DefaultInteractiveSelect.WithOptions(options).Show(question)
		case "multiselect":
			var options []string
			for _, option := range variable.OptionOrder {
				options = append(options, option)
			}
			value, err = pterm.DefaultInteractiveMultiselect.WithOptions(options).Show(question)
		case "component":
			// Components do not directly ask the user, only when they are used
		default:
			// Check if first letter is uppercase, if so, it is a component
			isComponent := variable.Type[0] >= 'A' && variable.Type[0] <= 'Z'
			if !isComponent {
				return fmt.Errorf("unknown variable type: %s", variable.Type)
			}

			fmt.Fprintf(os.Stderr, "Component: %s\n", variable.Type)

			// Find component by name
			var component *Variable
			for _, v := range variables {
				if v.Name == variable.Type {
					component = &v
					break
				}
			}

			// Ask for component variables
			var componentVariables []Variable
			for _, componentName := range *component.ComponentOrder {
				componentVariables = append(componentVariables, component.ComponentVars[componentName])
			}
			err = AskForVariables(componentVariables, variable.Name)
			if err != nil {
				return err
			}

			value = componentVariables
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
