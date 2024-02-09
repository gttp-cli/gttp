# Boolean

The `boolean` type can be used to define variables that accept boolean input.

## Basic

Basic syntax for the `boolean` type:

```yaml
variables:
  - name: AcceptedToS
    type: boolean # Set the type to boolean
    description: Do you accept the terms of service?
template: |-
  You have {{ if .AcceptedToS }}accepted{{ else }}not accepted{{ end }} the terms of service.
```
