---
sidebar_position: 1
---

# Template Syntax

GTTP uses YAML to define the structure of your template:

```yaml
template: Hello, World!
```

GTTP will parse this template to:

```
Hello, World!
```

## Variables

You can define variables in your template:

```yaml
variables:
  - name: Name
    type: text
    description: Name of the person
  - name: Lastname
    type: text
    description: Lastname of the person
```

You can use the defined variables in your template:

```yaml
variables:
  - name: Name
    type: text
    description: Name of the person
  - name: Lastname
    type: text
    description: Lastname of the person
template: Hello, {{ .Name }} {{ .Lastname }}!
```

When executing the template, GTTP will interactively ask you to fill out the defined variables:

```
Name of the person: John
Lastname of the person: Doe
```

When all variables are filled out, GTTP will parse the template to:

```
Hello, John Doe!
```

## Structures

Structures define custom data types.
They combine multiple variables into a single, reusable type.

```yaml
structures:
  person:
    - name: Name
      type: text
      description: Name of the person
    - name: Lastname
      type: text
      description: Lastname of the person
```

You can use the defined structures like any other type:

```yaml
structures:
  person:
    - name: Name
      type: text
      description: Name of the person
    - name: Lastname
      type: text
      description: Lastname of the person

variables:
    - name: UserA
      type: person
      description: A person
    - name: UserB
      type: person
      description: Another person

template: |-
    Hello, {{ .UserA.Name }} {{ .UserA.Lastname }}!
    Hello, {{ .UserB.Name }} {{ .UserB.Lastname }}!
```

When executing the template, GTTP will interactively ask you to fill out the defined variables:

```
[UserA] Name of the person: John
[UserA] Lastname of the person: Doe
[UserB] Name of the person: Jane
[UserB] Lastname of the person: Doe
```

When all variables are filled out, GTTP will parse the template to:
    
```
Hello, John Doe!
Hello, Jane Doe!
```
