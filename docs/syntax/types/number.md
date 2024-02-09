# Number

The `number` type can be used to define variables that accept numeric input.

## Basic

Basic syntax for the `number` type:

```yaml
variables:
  - name: Age
    type: number # Set the type to number
    description: Age of the person
template: |-
  You are {{ .Age }} years old.
```

## Validation

### Minimum

You can use the `min` property to define a minimum value for validation:

```yaml
variables:
  - name: Age
    type: number
    min: 18 # only allow ages 18 and above
    description: Age of the person
template: |-
  You are {{ .Age }} years old.
```

### Maximum

You can use the `max` property to define a maximum value for validation:

```yaml
variables:
  - name: Age
    type: number
    max: 99 # only allow ages 99 and below
    description: Age of the person
template: |-
  You are {{ .Age }} years old.
```
