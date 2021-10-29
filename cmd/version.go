package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the current CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
