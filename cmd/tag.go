package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/keygen-sh/keygen-cli/internal/keygenext"
	"github.com/spf13/cobra"
)

var (
	tagOpts = &TagCommandOptions{}
	tagCmd  = &cobra.Command{
		Use:   "tag <tag>",
		Short: "tag an existing release",
		Example: `  keygen tag latest \
      --account '1fddcec8-8dd3-4d8d-9b16-215cac0f9b52' \
      --product '2313b7e7-1ea6-4a01-901e-2931de6bb1e2' \
      --token 'prod-xxx' \
      --release '1.0.0'

Docs:
  https://keygen.sh/docs/cli/`,
		Args: tagArgs,
		RunE: tagRun,

		// Encountering an error should not display usage
		SilenceUsage: true,
	}
)

type TagCommandOptions struct {
	Release       string
	NoAutoUpgrade bool
	Package       string
}

func init() {
	tagCmd.Flags().StringVar(&keygenext.Account, "account", "", "your keygen.sh account identifier [$KEYGEN_ACCOUNT_ID=<id>] (required)")
	tagCmd.Flags().StringVar(&keygenext.Product, "product", "", "your keygen.sh product identifier [$KEYGEN_PRODUCT_ID=<id>] (required)")
	tagCmd.Flags().StringVar(&keygenext.Token, "token", "", "your keygen.sh product or environment token [$KEYGEN_TOKEN] (required)")
	tagCmd.Flags().StringVar(&keygenext.Environment, "environment", "", "your keygen.sh environment identifier [$KEYGEN_ENVIRONMENT=<id>]")
	tagCmd.Flags().StringVar(&keygenext.APIURL, "host", "", "the host of the keygen server [$KEYGEN_HOST=<host>]")
	tagCmd.Flags().StringVar(&tagOpts.Release, "release", "", "the release identifier (required)")
	tagCmd.Flags().StringVar(&tagOpts.Package, "package", "", "the package identifier")
	tagCmd.Flags().BoolVar(&tagOpts.NoAutoUpgrade, "no-auto-upgrade", false, "disable automatic upgrade checks [$KEYGEN_NO_AUTO_UPGRADE=1]")

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
		tagOpts.NoAutoUpgrade = true
	}

	if keygenext.Account == "" {
		tagCmd.MarkFlagRequired("account")
	}

	if keygenext.Product == "" {
		tagCmd.MarkFlagRequired("product")
	}

	if keygenext.Token == "" {
		tagCmd.MarkFlagRequired("token")
	}

	tagCmd.MarkFlagRequired("release")

	rootCmd.AddCommand(tagCmd)
}

func tagArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("tag is required")
	}

	return nil
}

func tagRun(cmd *cobra.Command, args []string) error {
	if !tagOpts.NoAutoUpgrade {
		err := upgradeRun(nil, nil)
		if err != nil {
			return err
		}
	}

	release := &keygenext.Release{
		ID:        tagOpts.Release,
		Tag:       &args[0],
		PackageID: &tagOpts.Package,
	}

	if err := release.Update(); err != nil {
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

	fmt.Println(green("tagged:") + " release " + italic(release.ID))

	return nil
}
