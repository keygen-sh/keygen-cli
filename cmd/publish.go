package cmd

import (
	"fmt"
	"os"

	"github.com/keygen-sh/keygen-cli/internal/keygenext"
	"github.com/spf13/cobra"
)

var (
	publishOpts = &PublishCommandOptions{}
	publishCmd  = &cobra.Command{
		Use:   "publish",
		Short: "publish an existing release",
		Example: `  keygen publish \
      --account '1fddcec8-8dd3-4d8d-9b16-215cac0f9b52' \
      --product '2313b7e7-1ea6-4a01-901e-2931de6bb1e2' \
      --token 'prod-xxx' \
      --release '1.0.0'

Docs:
  https://keygen.sh/docs/cli/`,
		Args: cobra.NoArgs,
		RunE: publishRun,

		// Encountering an error should not display usage
		SilenceUsage: true,
	}
)

type PublishCommandOptions struct {
	release       string
	noAutoUpgrade bool
}

func init() {
	publishCmd.Flags().StringVar(&keygenext.Account, "account", "", "your keygen.sh account identifier [$KEYGEN_ACCOUNT_ID=<id>] (required)")
	publishCmd.Flags().StringVar(&keygenext.Product, "product", "", "your keygen.sh product identifier [$KEYGEN_PRODUCT_ID=<id>] (required)")
	publishCmd.Flags().StringVar(&keygenext.Token, "token", "", "your keygen.sh product or environment token [$KEYGEN_TOKEN] (required)")
	publishCmd.Flags().StringVar(&keygenext.Environment, "environment", "", "your keygen.sh environment identifier [$KEYGEN_ENVIRONMENT=<id>]")
	publishCmd.Flags().StringVar(&publishOpts.release, "release", "", "the release identifier (required)")
	publishCmd.Flags().BoolVar(&publishOpts.noAutoUpgrade, "no-auto-upgrade", false, "disable automatic upgrade checks [$KEYGEN_NO_AUTO_UPGRADE=1]")

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

	if v, ok := os.LookupEnv("KEYGEN_ENVIRONMENT_TOKEN"); ok {
		if keygenext.Token == "" {
			keygenext.Token = v
		}
	}

	if v, ok := os.LookupEnv("KEYGEN_PRODUCT_TOKEN"); ok {
		if keygenext.Token == "" {
			keygenext.Token = v
		}
	}

	if v, ok := os.LookupEnv("KEYGEN_TOKEN"); ok {
		if keygenext.Token == "" {
			keygenext.Token = v
		}
	}

	if _, ok := os.LookupEnv("KEYGEN_NO_AUTO_UPGRADE"); ok {
		publishOpts.noAutoUpgrade = true
	}

	if keygenext.Account == "" {
		publishCmd.MarkFlagRequired("account")
	}

	if keygenext.Product == "" {
		publishCmd.MarkFlagRequired("product")
	}

	if keygenext.Token == "" {
		publishCmd.MarkFlagRequired("token")
	}

	publishCmd.MarkFlagRequired("release")

	rootCmd.AddCommand(publishCmd)
}

func publishRun(cmd *cobra.Command, args []string) error {
	if !publishOpts.noAutoUpgrade {
		err := upgradeRun(nil, nil)
		if err != nil {
			return err
		}
	}

	release := &keygenext.Release{
		ID: publishOpts.release,
	}

	if err := release.Publish(); err != nil {
		if e, ok := err.(*keygenext.Error); ok {
			var code string
			if e.Code != "" {
				code = italic("(" + e.Code + ")")
			}

			return fmt.Errorf("%s: %s %s", e.Title, e.Detail, code)
		}

		return err
	}

	fmt.Println(green("published:") + " release " + italic(release.ID))

	return nil
}
