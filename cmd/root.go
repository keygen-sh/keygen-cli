package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	// The current version of the program
	Version = "1.0.0"
)

var (
	rootCmd = &cobra.Command{
		Use:     "keygen",
		Short:   "CLI to interact with keygen.sh",
		Version: Version,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
