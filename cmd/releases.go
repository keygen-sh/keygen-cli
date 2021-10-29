package cmd

import (
	"github.com/spf13/cobra"
)

var (
	releasesCmd = &cobra.Command{
		Use:   "releases",
		Short: "Manage releases for your keygen.sh products",
	}
)

func init() {
	flags := releasesCmd.PersistentFlags()

	flags.StringP("account", "a", "", "Your keygen.sh account ID")
	flags.StringP("product", "p", "", "Your keygen.sh product ID")
	flags.StringP("token", "t", "", "Your keygen.sh product token")

	rootCmd.AddCommand(releasesCmd)
}
