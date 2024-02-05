<h1 align="center">💻 GTTP | Go Text Template Parser</h1>
<p align="center">A modern CLI to create and fill out reusable text templates</p>

<p align="center">

<a href="https://github.com/gttp-cli/gttp/releases" style="text-decoration: none">
    <img src="https://img.shields.io/github/v/release/gttp-cli/gttp?style=flat-square" alt="Latest Release">
</a>

<a href="https://github.com/gttp-cli/gttp/stargazers" style="text-decoration: none">
    <img src="https://img.shields.io/github/stars/gttp-cli/gttp.svg?style=flat-square" alt="Stars">
</a>

<a href="https://github.com/gttp-cli/gttp/fork" style="text-decoration: none">
    <img src="https://img.shields.io/github/forks/gttp-cli/gttp.svg?style=flat-square" alt="Forks">
</a>

<a href="https://opensource.org/licenses/MIT" style="text-decoration: none">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" alt="License: MIT">
</a>

<br/>

<a href="https://github.com/gttp-cli/gttp/releases" style="text-decoration: none">
    <img src="https://img.shields.io/badge/platform-windows%20%7C%20macos%20%7C%20linux-informational?style=for-the-badge" alt="Downloads">
</a>

 <a href="https://marvin.ws/twitter">
    <img src="https://img.shields.io/badge/Twitter-%40MarvinJWendt-1DA1F2?logo=twitter&style=for-the-badge"/>
</a>

<br/>
<br/>

</p>

## Introduction

GTTP lets you define your text templates using YAMl.  
When you execute a template file, GTTP will interactively ask you to fill out the defined variables.  
The template is then parsed with the Go [text/template](https://pkg.go.dev/text/template) syntax.

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

## Installation

There are multiple ways to install `gttp` on your system.

### Install Using Web Installer

You can install `gttp` using [instl](https://instl.sh).
Using instl is the simplest way to install `gttp` on your system.

Just copy the following command and paste it into your terminal:

| Platform | Command                                          |
|----------|--------------------------------------------------|
| Windows  | `iwr instl.sh/gttp-cli/gttp/windows \| iex`      |
| macOS    | `curl -sSL instl.sh/gttp-cli/gttp/macos \| bash` |
| Linux    | `curl -sSL instl.sh/gttp-cli/gttp/linux \| bash` |

> [!TIP]
> If you want to take a look at the script before running it, you can open the instl.sh URL in your browser.


### Install using Go

If you have [Go](https://go.dev) installed, you can install `gttp` using the following command:

```bash
go install github.com/gttp-cli/gttp@latest
```

## Docs

Docs are available at: https://docs.gttp.dev

