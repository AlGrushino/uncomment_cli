package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "uncomment-cli",
	Short: "is a powerful custom CLI tool",
	Long:  "A fast and flexible CLI tool built with Go and Cobra",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to uncomment-cli! Use --help to see available options.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute: %v", err)
	}
}
