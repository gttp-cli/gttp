# Multiselect

The `multiselect` type can be used to define variables that accept multiple options from a list of predefined options.

## Basic

Basic syntax for the `multiselect` type:

```yaml
variables:
  - name: Colors
    type: multiselect # Set the type to multiselect
    options: # Define the options
      - name: Red
      - name: Green
      - name: Blue
    description: Favorite colors
template: |-
  Your favorite colors are {{ .Colors }}.
```

## Custom values

You can use the `value` property to define custom values for the options:

```yaml
variables:
  - name: Colors
    type: multiselect
    options:
      - name: Red
        value: "#ff0000" # Set the value to a hex color code
      - name: Green
        value: "#00ff00"
      - name: Blue
        value: "#0000ff"
    description: Favorite colors
template: |-
  Your favorite colors are {{ .Colors }}.
```

## Iterate over results

You can use the `range` function in the template to iterate over the results:

```yaml
variables:
  - name: Colors
    type: multiselect
    options:
      - name: Red
      - name: Green
      - name: Blue
    description: Favorite colors
template: |-
  Your favorite colors are:
  {{- range .Colors }}
  - {{ . }}
  {{- end }}
```

When every color is selected, the output will be:

```
Your favorite colors are:
- Red
- Green
- Blue
```
