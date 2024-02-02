package model

import (
	"encoding/json"
	"github.com/goccy/go-yaml"
)

type Template struct {
	Structures map[string][]Variable `json:"structures,omitempty"`
	Variables  []Variable            `json:"variables,omitempty"`

	Template string `json:"template"`
}

type Variable struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	IsArray     bool   `json:"array,omitempty"`
	Description string `json:"description,omitempty"`

	Value   any `json:"value,omitempty"`
	Default any `json:"default,omitempty"`

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
