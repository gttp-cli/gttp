# Text

The `text` type can be used to define variables that accept text input.

## Basic

Basic syntax for the `text` type:

```yaml
variables:
  - name: Name
    type: text # Set the type to text
    description: Name of the person
template: |-
  Hello, {{ .Name }}!
```

## Multiline

You can use the `multiline` property to define a multiline text input:

```yaml
variables:
  - name: Description
    type: text
    multiline: true # Set the multiline property to true
    description: Description of the person
template: |-
  Description:
  {{ .Description }}
```

## Validation

### Regex

You can use the `regex` property to define a regular expression for validation:

```yaml
variables:
  - name: Text
    type: text
    regex: ^[a-z]+$ # only allow lowercase letters
    description: A string of lowercase letters
template: |-
  {{ .Text }}
```
