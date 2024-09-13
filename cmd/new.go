package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/keygen-sh/keygen-cli/internal/keygenext"
	"github.com/spf13/cobra"
)

var (
	draftOpts = &DraftCommandOptions{}
	draftCmd  = &cobra.Command{
		Use:   "new",
		Short: "create a new draft release",
		Example: `  keygen new \
      --account '1fddcec8-8dd3-4d8d-9b16-215cac0f9b52' \
      --product '2313b7e7-1ea6-4a01-901e-2931de6bb1e2' \
      --token 'prod-xxx' \
      --channel 'beta' \
      --version '1.0.0'

Docs:
  https://keygen.sh/docs/cli/`,
		Args: cobra.NoArgs,
		RunE: draftRun,

		// Encountering an error should not display usage
		SilenceUsage: true,
	}
)

type DraftCommandOptions struct {
	Name          string
	Description   string
	Version       string
	Tag           string
	Channel       string
	Package       string
	Entitlements  []string
	NoAutoUpgrade bool
	Metadata      map[string]string // Add metadata field
}

func init() {
	draftCmd.Flags().StringVar(&keygenext.Account, "account", "", "your keygen.sh account identifier [$KEYGEN_ACCOUNT_ID=<id>] (required)")
	draftCmd.Flags().StringVar(&keygenext.Product, "product", "", "your keygen.sh product identifier [$KEYGEN_PRODUCT_ID=<id>] (required)")
	draftCmd.Flags().StringVar(&keygenext.Token, "token", "", "your keygen.sh product or environment token [$KEYGEN_TOKEN] (required)")
	draftCmd.Flags().StringVar(&keygenext.Environment, "environment", "", "your keygen.sh environment identifier [$KEYGEN_ENVIRONMENT=<id>]")
	draftCmd.Flags().StringVar(&keygenext.APIURL, "host", "", "the host of the keygen server [$KEYGEN_HOST=<host>]")
	draftCmd.Flags().StringVar(&draftOpts.Version, "version", "", "version for the release (required)")
	draftCmd.Flags().StringVar(&draftOpts.Tag, "tag", "", "tag for the release")
	draftCmd.Flags().StringVar(&draftOpts.Name, "name", "", "human-readable name for the release")
	draftCmd.Flags().StringVar(&draftOpts.Description, "description", "", "description for the release (e.g. release notes)")
	draftCmd.Flags().StringVar(&draftOpts.Channel, "channel", "stable", "channel for the release, one of: stable, rc, beta, alpha, dev")
	draftCmd.Flags().StringVar(&draftOpts.Package, "package", "", "package identifier for the release")
	draftCmd.Flags().BoolVar(&draftOpts.NoAutoUpgrade, "no-auto-upgrade", false, "disable automatic upgrade checks [$KEYGEN_NO_AUTO_UPGRADE=1]")
	draftCmd.Flags().StringSliceVar(&draftOpts.Entitlements, "entitlements", []string{}, "comma seperated list of entitlement constraints (e.g. --entitlements <id>,<id>,...)")
	draftCmd.Flags().StringToStringVar(&draftOpts.Metadata, "metadata", map[string]string{}, "metadata for the release as key-value pairs (e.g. --metadata key1=value1,key2=value2)")

	// TODO(ezekg) Prompt multi-line description input from stdin if "--"?
	// TODO(ezekg) Add metadata flag

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
		draftOpts.NoAutoUpgrade = true
	}

	if keygenext.Account == "" {
		draftCmd.MarkFlagRequired("account")
	}

	if keygenext.Product == "" {
		draftCmd.MarkFlagRequired("product")
	}

	if keygenext.Token == "" {
		draftCmd.MarkFlagRequired("token")
	}

	draftCmd.MarkFlagRequired("version")

	rootCmd.AddCommand(draftCmd)
}

func draftRun(cmd *cobra.Command, args []string) error {
	if !draftOpts.NoAutoUpgrade {
		err := upgradeRun(nil, nil)
		if err != nil {
			return err
		}
	}

	channel := draftOpts.Channel

	var constraints keygenext.Constraints
	if e := draftOpts.Entitlements; len(e) != 0 {
		constraints = constraints.From(e)
	}

	var tag *string
	if t := draftOpts.Tag; t != "" {
		tag = &t
	}

	var name *string
	if n := draftOpts.Name; n != "" {
		name = &n
	}

	var desc *string
	if d := draftOpts.Description; d != "" {
		desc = &d
	}

	var pkg *string
	if p := draftOpts.Package; p != "" {
		pkg = &p
	}

	version, err := semver.NewVersion(draftOpts.Version)
	if err != nil {
		return fmt.Errorf(`version "%s" is not acceptable (%s)`, draftOpts.Version, italic(strings.ToLower(err.Error())))
	}

	release := &keygenext.Release{
		Name:        name,
		Description: desc,
		Version:     version.String(),
		Tag:         tag,
		Channel:     channel,
		ProductID:   keygenext.Product,
		PackageID:   pkg,
		Constraints: constraints,
		Metadata:    draftOpts.Metadata, // Include the metadata here
	}

	if err := release.Create(); err != nil {
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

	fmt.Println(green("drafted:") + " release " + italic(release.ID))

	return nil
}
