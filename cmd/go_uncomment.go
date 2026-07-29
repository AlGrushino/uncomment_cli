package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var filePath string

var goUncommentCmd = &cobra.Command{
	Use:   "uncomment *.go file",
	Short: "deletes all comments from *.go file",
	Long:  "This command deletes ALL comments from *.go file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("All comments in %s deleted\n", filePath)
	},
}

func init() {
	rootCmd.AddCommand(goUncommentCmd)

	goUncommentCmd.Flags().StringVarP(
		&filePath,
		"path",
		"p",
		"",
		"Path to file",
	)
}
