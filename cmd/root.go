package cmd

import (
	"os"

	"github.com/keygen-sh/keygen-cli/internal/keygenext"
	"github.com/spf13/cobra"
)

const (
	// The current version of the CLI.
	Version = "1.0.0-beta.1"
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

type CommandOptions struct {
	filename         string
	name             string
	description      string
	version          string
	platform         string
	channel          string
	entitlements     []string
	signature        string
	checksum         string
	signingAlgorithm string
	signingKey       string
	verifyKey        string
}

func init() {
	keygenext.UserAgent = "cli/" + Version
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
