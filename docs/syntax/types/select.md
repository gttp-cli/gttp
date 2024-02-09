# Select

The `select` type can be used to define variables that accept a single option from a list of predefined options.

## Basic

Basic syntax for the `select` type:

```yaml
variables:
  - name: Color
    type: select # Set the type to select
    options: # Define the options
      - name: Red
      - name: Green
      - name: Blue
    description: Favorite color
template: |-
  Your favorite color is {{ .Color }}.
```

## Custom values

You can use the `value` property to define custom values for the options:

```yaml
variables:
  - name: Color
    type: select
    options:
      - name: Red
        value: "#ff0000" # Set the value to a hex color code
      - name: Green
        value: "#00ff00"
      - name: Blue
        value: "#0000ff"
    description: Favorite color
template: |-
  Your favorite color is {{ .Color }}.
```
