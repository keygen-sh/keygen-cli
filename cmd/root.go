package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/keygen-sh/keygen-cli/internal/keygenext"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "keygen",
		Short: "CLI to interact with keygen.sh",
		Long: `CLI to interact with keygen.sh

Version:
  keygen/` + Version + " " + runtime.GOOS + "-" + runtime.GOARCH + " " + runtime.Version(),
		Version:       Version,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
)

func init() {
	keygenext.UserAgent = "cli/" + Version

	rootCmd.PersistentFlags().BoolVar(&color.NoColor, "no-color", false, "disable colors in command output [$NO_COLOR=1]")

	rootCmd.InitDefaultVersionFlag()
	rootCmd.InitDefaultHelpFlag()

	rootCmd.SetHelpCommand(helpCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, red("error:")+" "+err.Error())

		os.Exit(1)
	}
}
