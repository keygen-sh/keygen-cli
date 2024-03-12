package cmd

import (
	"fmt"
	"os"

	"github.com/keygen-sh/keygen-cli/internal/keygenext"
	"github.com/spf13/cobra"
)

var (
	delOpts = &DeleteCommandOptions{}
	delCmd  = &cobra.Command{
		Use:   "del",
		Short: "delete an existing release",
		Example: `  keygen del \
      --account '1fddcec8-8dd3-4d8d-9b16-215cac0f9b52' \
      --product '2313b7e7-1ea6-4a01-901e-2931de6bb1e2' \
      --token 'prod-xxx' \
      --release '1.0.0'

Docs:
  https://keygen.sh/docs/cli/`,
		Args: cobra.NoArgs,
		RunE: delRun,

		// Encountering an error should not display usage
		SilenceUsage: true,
	}
)

type DeleteCommandOptions struct {
	Release       string
	NoAutoUpgrade bool
}

func init() {
	delCmd.Flags().StringVar(&keygenext.Account, "account", "", "your keygen.sh account identifier [$KEYGEN_ACCOUNT_ID=<id>] (required)")
	delCmd.Flags().StringVar(&keygenext.Product, "product", "", "your keygen.sh product identifier [$KEYGEN_PRODUCT_ID=<id>] (required)")
	delCmd.Flags().StringVar(&keygenext.Token, "token", "", "your keygen.sh environment or product token [$KEYGEN_TOKEN] (required)")
	delCmd.Flags().StringVar(&keygenext.Environment, "environment", "", "your keygen.sh environment identifier [$KEYGEN_ENVIRONMENT=<id>]")
	delCmd.Flags().StringVar(&keygenext.APIURL, "host", "", "the host of the keygen server [$KEYGEN_HOST=<host>]")
	delCmd.Flags().StringVar(&delOpts.Release, "release", "", "the release identifier (required)")
	delCmd.Flags().BoolVar(&delOpts.NoAutoUpgrade, "no-auto-upgrade", false, "disable automatic upgrade checks [$KEYGEN_NO_AUTO_UPGRADE=1]")

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

	if v, ok := os.LookupEnv("KEYGEN_HOST"); ok {
		if keygenext.APIURL == "" {
			keygenext.APIURL = v
		}
	}

	if _, ok := os.LookupEnv("KEYGEN_NO_AUTO_UPGRADE"); ok {
		delOpts.NoAutoUpgrade = true
	}

	if keygenext.Account == "" {
		delCmd.MarkFlagRequired("account")
	}

	if keygenext.Product == "" {
		delCmd.MarkFlagRequired("product")
	}

	if keygenext.Token == "" {
		delCmd.MarkFlagRequired("token")
	}

	delCmd.MarkFlagRequired("release")

	rootCmd.AddCommand(delCmd)
}

func delRun(cmd *cobra.Command, args []string) error {
	if !delOpts.NoAutoUpgrade {
		err := upgradeRun(nil, nil)
		if err != nil {
			return err
		}
	}

	release := &keygenext.Release{
		ID: delOpts.Release,
	}

	if err := release.Delete(); err != nil {
		if e, ok := err.(*keygenext.Error); ok {
			var code string
			if e.Code != "" {
				code = italic("(" + e.Code + ")")
			}

			if e.Source != "" {
				return fmt.Errorf("%s: %s %s %s", e.Title, e.Source, e.Detail, code)
			} else {
				return fmt.Errorf("%s: %s %s", e.Title, e.Detail, code)
			}
		}

		return err
	}

	fmt.Println(green("deleted:") + " release " + italic(release.ID))

	return nil
}
