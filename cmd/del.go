package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/keygen-sh/jsonapi-go"
	"github.com/keygen-sh/keygen-cli/internal/keygenext"
	"github.com/spf13/cobra"
)

var (
	delOpts = &DeleteCommandOptions{}
	delCmd  = &cobra.Command{
		Use:   "del",
		Short: "delete an existing release or artifact",
		Example: `  keygen del \
      --account '1fddcec8-8dd3-4d8d-9b16-215cac0f9b52' \
      --product '2313b7e7-1ea6-4a01-901e-2931de6bb1e2' \
      --token 'prod-xxx' \
      --release '1.0.0'

Docs:
  https://keygen.sh/docs/cli/`,
		Args:         cobra.NoArgs,
		RunE:         delRun,
		SilenceUsage: true,
	}
)

type DeleteCommandOptions struct {
	Release       string
	Package       string
	Artifact      string
	NoAutoUpgrade bool
}

func init() {
	delCmd.Flags().StringVar(&keygenext.Account, "account", "", "your keygen.sh account identifier [$KEYGEN_ACCOUNT_ID=<id>] (required)")
	delCmd.Flags().StringVar(&keygenext.Product, "product", "", "your keygen.sh product identifier [$KEYGEN_PRODUCT_ID=<id>] (required)")
	delCmd.Flags().StringVar(&keygenext.Token, "token", "", "your keygen.sh environment or product token [$KEYGEN_TOKEN] (required)")
	delCmd.Flags().StringVar(&keygenext.Environment, "environment", "", "your keygen.sh environment identifier [$KEYGEN_ENVIRONMENT=<id>]")
	delCmd.Flags().StringVar(&keygenext.APIURL, "host", "", "the host of the keygen server [$KEYGEN_HOST=<host>]")
	delCmd.Flags().StringVar(&delOpts.Release, "release", "", "the release identifier (required)")
	delCmd.Flags().StringVar(&delOpts.Package, "package", "", "package identifier for the release")
	delCmd.Flags().StringVar(&delOpts.Artifact, "artifact", "", "the artifact identifier")
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

	var deletable interface {
		jsonapi.MarshalResourceIdentifier
		Delete() error
	}

	switch {
	case delOpts.Artifact != "":
		deletable = &keygenext.Artifact{ReleaseID: &delOpts.Release, ID: delOpts.Artifact}
	default:
		release := &keygenext.Release{
			ID:        delOpts.Release,
			PackageID: &delOpts.Package,
		}

		if err := release.Get(); err != nil {
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

		deletable = release
	}

	if err := deletable.Delete(); err != nil {
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

	fmt.Println(
		green("deleted:") + " " + strings.TrimSuffix(deletable.GetType(), "s") + " " + italic(deletable.GetID()),
	)

	return nil
}
