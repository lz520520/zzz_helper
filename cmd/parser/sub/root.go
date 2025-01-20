package sub

import (
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "parser",
	Short: "parser.",

	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {

}
