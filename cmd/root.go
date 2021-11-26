package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/keygen-sh/keygen-cli/internal/keygenext"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:           "keygen",
		Short:         "CLI to interact with keygen.sh",
		Version:       Version,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
)

type CommandOptions struct {
	filename         string
	filetype         string
	name             string
	description      string
	version          string
	platform         string
	channel          string
	entitlements     []string
	signature        string
	checksum         string
	signingAlgorithm string
	signingKeyPath   string
	verifyKeyPath    string
	signingKey       string
	noAutoUpgrade    bool
}

func init() {
	keygenext.UserAgent = "cli/" + Version

	rootCmd.PersistentFlags().BoolVar(&color.NoColor, "no-color", false, "disable colors in command output [$NO_COLOR=1]")

	rootCmd.InitDefaultVersionFlag()
	rootCmd.InitDefaultHelpFlag()

	rootCmd.SetHelpCommand(helpCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		red := color.New(color.FgRed).SprintFunc()

		fmt.Fprintln(os.Stderr, red("error:")+" "+err.Error())

		os.Exit(1)
	}
}
