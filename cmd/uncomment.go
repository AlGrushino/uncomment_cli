package cmd

import (
	"uncomment-cli/internal/service"

	"github.com/spf13/cobra"
)

var filePath string

var uncommentCmd = &cobra.Command{
	Use:   "uncomment",
	Short: "deletes all comments from *.go file",
	Long:  "This command deletes ALL comments from *.go file",

	Run: func(cmd *cobra.Command, args []string) {
		err := service.Uncomment(filePath)
		if err != nil {
			cmd.PrintErrln("Error:", err)
			return
		}
		cmd.Printf("All comments in %s deleted\n", filePath)
	},
}

func init() {
	rootCmd.AddCommand(uncommentCmd)

	uncommentCmd.Flags().StringVarP(
		&filePath,
		"path",
		"p",
		"",
		"Path to file",
	)
}
