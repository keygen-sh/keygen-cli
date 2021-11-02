package cmd

import (
	"github.com/keygen-sh/cli/keygenext"
	"github.com/spf13/cobra"
)

var (
	releasesCmd = &cobra.Command{
		Use:   "releases",
		Short: "Manage releases for your keygen.sh products",
	}
)

func init() {
	releasesCmd.PersistentFlags().StringVar(&keygenext.Account, "account", "", "your keygen.sh account identifier (required)")
	releasesCmd.PersistentFlags().StringVar(&keygenext.Product, "product", "", "your keygen.sh product identifier (required)")
	releasesCmd.PersistentFlags().StringVar(&keygenext.Token, "token", "", "your keygen.sh product token (required)")

	releasesCmd.MarkPersistentFlagRequired("account")
	releasesCmd.MarkPersistentFlagRequired("product")
	releasesCmd.MarkPersistentFlagRequired("token")

	rootCmd.AddCommand(releasesCmd)
}
