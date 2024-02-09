package model

import (
	"fmt"
	"regexp"
)

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
		// Default vaue must be an float or int or nil
		if v.Default != nil {
			switch v.Default.(type) {
			case float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				// noop
			default:
				errors = append(errors, newValidationError(v, fmt.Sprintf("default must be a number or nil, got %T", v.Default)))
			}
		}

		var value float64
		if v.Value != nil {
			var ok bool
			value, ok = v.Value.(float64)
			if !ok {
				errors = append(errors, newValidationError(v, fmt.Sprintf("value must be a number, got %T", v.Value)))
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
		// Default value must be a string or nil
		if v.Default != nil {
			_, ok := v.Default.(string)
			if !ok {
				errors = append(errors, newValidationError(v, fmt.Sprintf("default must be a string or nil, got %T", v.Default)))
			}
		}

		var value string
		if v.Value != nil {
			var ok bool
			value, ok = v.Value.(string)
			if !ok {
				errors = append(errors, newValidationError(v, fmt.Sprintf("value must be a string, got %T", v.Value)))
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
