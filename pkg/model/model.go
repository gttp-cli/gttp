package model

import (
	"encoding/json"
	"fmt"
	"github.com/goccy/go-yaml"
	"regexp"
)

type Template struct {
	Structures map[string][]Variable `json:"structures,omitempty"`
	Variables  []Variable            `json:"variables,omitempty" jsonschema:"required"`

	Template string `json:"template" jsonschema:"required"`
}

type Variable struct {
	Name        string `json:"name" jsonschema:"required"`
	Type        string `json:"type" jsonschema:"required"`
	IsArray     bool   `json:"array,omitempty"`
	Description string `json:"description,omitempty"`

	Condition string `json:"condition,omitempty"`

	Value   any `json:"value,omitempty"`
	Default any `json:"default,omitempty"`

	// Number validation
	Min float64 `json:"min,omitempty"`
	Max float64 `json:"max,omitempty"`

	// String validation
	Regex string `json:"regex,omitempty"`

	Options []Option `json:"options,omitempty"`
}

type Option struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	Value   any `json:"value,omitempty"`
	Default any `json:"default,omitempty"`
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
