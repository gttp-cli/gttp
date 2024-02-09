package model

import (
	"encoding/json"
	"github.com/goccy/go-yaml"
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
