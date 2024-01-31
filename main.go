package main

import (
	"github.com/gttp-cli/gttp/cmd"
	"github.com/pterm/pterm"
)

func main() {
	pterm.Error.PrintOnError(cmd.Execute())
}
