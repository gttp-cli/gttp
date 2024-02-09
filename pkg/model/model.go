package model

import (
	"encoding/json"
	"fmt"
	"github.com/goccy/go-yaml"
	"regexp"
)

type Template struct {
	// Structures define custom types.
	// They can be used as reusable components, consiting of multiple variables.
	Structures map[string][]Variable `json:"structures,omitempty"`

	// Variables define the input variables for the template.
	// They can be used in the template content and in conditions.
	Variables []Variable `json:"variables"`

	// Template defines the content of the template.
	Template string `json:"template"`
}

type Variable struct {
	// Name is the name of the variable.
	Name string `json:"name"`
	// Type is the type of the variable.
	Type string `json:"type"`
	// IsArray indicates if the variable is an array.
	// Can also be indicated by the type, e.g. "string[]".
	IsArray bool `json:"array,omitempty"`
	// Multiline indicates if the variable is a multiline string.
	// Only applicable to text types.
	Multiline bool `json:"multiline,omitempty"`
	// Description is a description of the variable.
	Description string `json:"description,omitempty"`

	// Condition is a condition that must be met for the variable to be used.
	// Conditions are evaluated using expr-lang expressions (see: https://expr-lang.org/).
	Condition string `json:"condition,omitempty"`

	// Value is the value of the variable.
	// If the value is predefined in the template, the user will not be asked for input.
	Value any `json:"value,omitempty"`
	// Default is the default value of the variable, if the user does not provide a value.
	Default any `json:"default,omitempty"`

	// Min is the minimum value of the variable.
	// Only applicable to number types.
	Min float64 `json:"min,omitempty"`
	// Max is the maximum value of the variable.
	// Only applicable to number types.
	Max float64 `json:"max,omitempty"`

	// Regex is a regular expression that the value must match.
	// Only applicable to text types.
	Regex string `json:"regex,omitempty"`

	// Options are the available options for select and multiselect types.
	Options []Option `json:"options,omitempty"`
}

type Option struct {
	// Name is the name of the option.
	// If no value is provided, the name will be used as the value.
	Name string `json:"name"`

	// Value is the value of the option.
	// If no value is provided, the name will be used as the value.
	Value any `json:"value,omitempty"`
}

func (t Template) ToJSON() (string, error) {
	j, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return "", err
	}

	return string(j), nil
}

func (t Template) ToYAML() (string, error) {
	y, err := yaml.Marshal(t)
	if err != nil {
		return "", err
	}

	return string(y), nil
}

func FromJSON(jsonString string) (Template, error) {
	var t Template
	err := json.Unmarshal([]byte(jsonString), &t)
	if err != nil {
		return Template{}, err
	}

	return t, nil
}

func FromYAML(yamlString string) (Template, error) {
	var t Template
	err := yaml.Unmarshal([]byte(yamlString), &t)
	if err != nil {
		return Template{}, err
	}

	return t, nil
}

func (v Variable) Validate() []error {
	var errors []error

	// Check that type is set
	if v.Type == "" {
		errors = append(errors, newValidationError(v, "type is required"))
	}

	// Min and max are only applicable to number types
	if v.Min != 0 || v.Max != 0 {
		if v.Type != "number" {
			errors = append(errors, newValidationError(v, "min and max are only applicable to number types"))
		}
	}

	// Regex is only applicable to text types
	if v.Regex != "" {
		if v.Type != "text" {
			errors = append(errors, newValidationError(v, "regex is only applicable to text types"))
		}
	}

	// Multiline is only applicable to text types
	if v.Multiline {
		if v.Type != "text" {
			errors = append(errors, newValidationError(v, "multiline is only applicable to text types"))
		}
	}

	// Options are only applicable to select and multiselect types
	if len(v.Options) > 0 {
		if v.Type != "select" && v.Type != "multiselect" {
			errors = append(errors, newValidationError(v, "options are only applicable to select and multiselect types"))
		}
	}

	// Specific type validations
	switch v.Type {
	case "number":
		var value float64
		if v.Value != nil {
			var ok bool
			value, ok = v.Value.(float64)
			if !ok {
				errors = append(errors, newValidationError(v, "value must be a number"))
			}
		}

		if v.Min != 0 || v.Max != 0 { // min or max is set
			if v.Min > v.Max {
				errors = append(errors, newValidationError(v, "min must be less than max"))
			}

			if v.Value != nil {
				if value < v.Min || value > v.Max {
					errors = append(errors, newValidationError(v, "value must be between min and max"))
				}
			}
		}
	case "text":
		var value string
		if v.Value != nil {
			var ok bool
			value, ok = v.Value.(string)
			if !ok {
				errors = append(errors, newValidationError(v, "value must be a string"))
			}
		}

		// Validate regex
		if v.Regex != "" {
			if v.Value != nil {
				re, err := regexp.Compile(v.Regex)
				if err != nil {
					errors = append(errors, newValidationError(v, "invalid regex"))
				}

				if !re.MatchString(value) {
					errors = append(errors, newValidationError(v, "value does not match regex"))
				}
			}
		}

	case "boolean":
		if v.Value != nil {
			_, ok := v.Value.(bool)
			if !ok {
				errors = append(errors, newValidationError(v, "value must be a boolean"))
			}
		}

	case "select":
		if len(v.Options) == 0 {
			errors = append(errors, newValidationError(v, "options are required"))
		}

		if v.Value != nil {
			var ok bool
			_, ok = v.Value.(string)
			if !ok {
				errors = append(errors, newValidationError(v, "value must be a string"))
			}
		}

	case "multiselect":
		if len(v.Options) == 0 {
			errors = append(errors, newValidationError(v, "options are required"))
		}

		if v.Value != nil {
			// Value must be either string or string slice
			switch v.Value.(type) {
			case string:
				// Check if value is in options
				found := false
				for _, o := range v.Options {
					if v.Value == o.Value {
						found = true
						break
					}
				}

				if !found {
					errors = append(errors, newValidationError(v, "value is not in options"))
				}

			case []string:
				// Check if all values are in options
				for _, value := range v.Value.([]string) {
					found := false
					for _, o := range v.Options {
						if value == o.Value {
							found = true
							break
						}
					}

					if !found {
						errors = append(errors, newValidationError(v, "value is not in options"))
					}
				}

			default:
				errors = append(errors, newValidationError(v, "value must be a string or string slice"))
			}
		}

	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (t Template) Validate() []error {
	var errors []error

	if t.Template == "" {
		errors = append(errors, fmt.Errorf("template is required"))
	}

	for _, v := range t.Variables {
		errs := v.Validate()
		if errs != nil {
			errors = append(errors, errs...)
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func newValidationError(v Variable, message string) error {
	return fmt.Errorf("variable %s: %s", v.Name, message)
}
