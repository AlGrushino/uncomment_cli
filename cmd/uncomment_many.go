package cmd

import (
	"uncomment-cli/internal/service"

	"github.com/spf13/cobra"
)

var pathList []string

var uncommentManyCmd = &cobra.Command{
	Use:   "uncomment-many",
	Short: "deletes all comments from any number of *.go files",
	Long:  "This command deletes ALL comments from any number of *.go files",

	Run: func(cmd *cobra.Command, args []string) {
		errCh := service.UncommentMany(pathList...)

		for p := range errCh {
			if p.Err != nil {
				cmd.PrintErrf("Failed to uncomment file: %s, error: %v\n", p.Path, p.Err)
				continue
			}
			cmd.Printf("Succesfully uncommented file: %s\n", p.Path)
		}
	},
}

func init() {
	rootCmd.AddCommand(uncommentManyCmd)

	uncommentManyCmd.Flags().StringArrayVarP(
		&pathList,
		"pathes",
		"p",
		[]string{},
		"Pathes to files",
	)
}
