package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// The current version of the CLI, embedded at compile time.
	Version string
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print the current CLI version",
		Args:  cobra.NoArgs,
		Run:   versionRun,
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func versionRun(cmd *cobra.Command, args []string) {
	fmt.Println(Version)
}
