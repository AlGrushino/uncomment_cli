package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var userLabel string

var greetCmd = &cobra.Command{
	Use:   "greet",
	Short: "Prints a friendly greeting",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Hello, %s!\n", userLabel)
	},
}

func init() {
	rootCmd.AddCommand(greetCmd)

	greetCmd.Flags().StringVarP(
		&userLabel,
		"name",
		"n",
		"World",
		"Name of the person to greet",
	)
}
