package main

import (
	"uncomment-cli/cmd"
)

// main is the entry point of uncomment-cli; it delegates execution to the
// root Cobra command.
func main() {
	cmd.Execute()
}
