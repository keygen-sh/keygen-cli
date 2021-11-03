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
	options = &CommandOptions{}
	rootCmd = &cobra.Command{
		Use:     "keygen",
		Short:   "CLI to interact with keygen.sh",
		Version: Version,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
)

type CommandOptions struct {
	filename    string
	name        string
	version     string
	platform    string
	channel     string
	constraints []string
	privateKey  string
	publicKey   string
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
