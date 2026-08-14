package main

import (
	"github.com/AlGrushino/uncomment_cli/cmd"
)

// main is the entry point of uncomment-cli; it delegates execution to the
// root Cobra command.
func main() {
	cmd.Execute()
}
