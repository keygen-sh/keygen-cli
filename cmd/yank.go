package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/keygen-sh/keygen-cli/internal/keygenext"
	"github.com/spf13/cobra"
)

var (
	yankOpts = &CommandOptions{}
	yankCmd  = &cobra.Command{
		Use:   "yank",
		Short: "yank an existing release",
		Example: `  keygen yank \
      --account '1fddcec8-8dd3-4d8d-9b16-215cac0f9b52' \
      --product '2313b7e7-1ea6-4a01-901e-2931de6bb1e2' \
      --token 'prod-xxx' \
      --release '1.0.0'

Docs:
  https://keygen.sh/docs/cli/`,
		Args: cobra.NoArgs,
		RunE: yankRun,

		// Encountering an error should not display usage
		SilenceUsage: true,
	}
)

func init() {
	yankCmd.Flags().StringVar(&keygenext.Account, "account", "", "your keygen.sh account identifier [$KEYGEN_ACCOUNT_ID=<id>] (required)")
	yankCmd.Flags().StringVar(&keygenext.Product, "product", "", "your keygen.sh product identifier [$KEYGEN_PRODUCT_ID=<id>] (required)")
	yankCmd.Flags().StringVar(&keygenext.Token, "token", "", "your keygen.sh product token [$KEYGEN_PRODUCT_TOKEN] (required)")
	yankCmd.Flags().StringVar(&yankOpts.release, "release", "", "the release identifier (required)")
	yankCmd.Flags().BoolVar(&yankOpts.noAutoUpgrade, "no-auto-upgrade", false, "disable automatic upgrade checks [$KEYGEN_NO_AUTO_UPGRADE=1]")

	if v, ok := os.LookupEnv("KEYGEN_ACCOUNT_ID"); ok {
		if keygenext.Account == "" {
			keygenext.Account = v
		}
	}

	if v, ok := os.LookupEnv("KEYGEN_PRODUCT_ID"); ok {
		if keygenext.Product == "" {
			keygenext.Product = v
		}
	}

	if v, ok := os.LookupEnv("KEYGEN_PRODUCT_TOKEN"); ok {
		if keygenext.Token == "" {
			keygenext.Token = v
		}
	}

	if _, ok := os.LookupEnv("KEYGEN_NO_AUTO_UPGRADE"); ok {
		yankOpts.noAutoUpgrade = true
	}

	if keygenext.Account == "" {
		yankCmd.MarkFlagRequired("account")
	}

	if keygenext.Product == "" {
		yankCmd.MarkFlagRequired("product")
	}

	if keygenext.Token == "" {
		yankCmd.MarkFlagRequired("token")
	}

	yankCmd.MarkFlagRequired("release")

	rootCmd.AddCommand(yankCmd)
}

func yankRun(cmd *cobra.Command, args []string) error {
	if !yankOpts.noAutoUpgrade {
		err := upgradeRun(nil, nil)
		if err != nil {
			return err
		}
	}

	italic := color.New(color.Italic).SprintFunc()
	release := &keygenext.Release{
		ID: yankOpts.release,
	}

	if err := release.Yank(); err != nil {
		e, ok := err.(*keygenext.Error)
		if ok {
			var code string
			if e.Code != "" {
				code = italic("(" + e.Code + ")")
			}

			return fmt.Errorf("%s: %s %s", e.Title, e.Detail, code)
		}

		return err
	}

	fmt.Println("yanked release " + italic(release.ID))

	return nil
}
