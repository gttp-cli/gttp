---
sidebar_position: 2
---

# 🚀 Run Demo Template

```yaml
# yaml-language-server: $schema=https://gttp.dev/schema

structures:
  person:
    - name: Name
      type: text
      description: Name of the person
    - name: Admin
      type: boolean
      description: Is the person an admin

variables:
  - name: Users
    type: person[]

template: |-
  You have added the following users:
  {{ range .Users }}
  - {{ .Name }} is an admin: {{ .Admin }}
  {{ end }}
```

This is a simple example of a GTTP template.  
When you execute a template file, GTTP will interactively ask you to fill out the defined variables.

**Run Demo via GTTP CLI:**

:::note
You need to have GTTP installed on your system to run the following command.

If you don't have GTTP installed, you can use Docker to run the demo.
:::

```bash
gttp -u gttp.dev/demo.yml
```

**Run Demo via Docker:**

```bash
docker run -it --rm ghcr.io/gttp-cli/gttp:main -u gttp.dev/demo.yml
```
