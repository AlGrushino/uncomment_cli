// Package cmd implements the CLI commands of uncomment-cli using Cobra:
// the root command, uncomment, uncomment-many and greet.
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// version is the CLI version. It can be overridden at build time via
// -ldflags "-X github.com/AlGrushino/uncomment_cli/cmd.version=<version>".
var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:     "uncomment-cli",
	Short:   "is a powerful custom CLI tool",
	Long:    "A fast and flexible CLI tool built with Go and Cobra",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to uncomment-cli! Use --help to see available options.")
	},
}

// Execute runs the root command. If it fails, the error is logged and the
// process terminates with a fatal status.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute: %v", err)
	}
}
