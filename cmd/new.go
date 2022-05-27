package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"
	"github.com/keygen-sh/keygen-cli/internal/keygenext"
	"github.com/spf13/cobra"
)

var (
	draftOpts = &CommandOptions{}
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

func init() {
	draftCmd.Flags().StringVar(&keygenext.Account, "account", "", "your keygen.sh account identifier [$KEYGEN_ACCOUNT_ID=<id>] (required)")
	draftCmd.Flags().StringVar(&keygenext.Product, "product", "", "your keygen.sh product identifier [$KEYGEN_PRODUCT_ID=<id>] (required)")
	draftCmd.Flags().StringVar(&keygenext.Token, "token", "", "your keygen.sh product token [$KEYGEN_PRODUCT_TOKEN] (required)")
	draftCmd.Flags().StringVar(&draftOpts.version, "version", "", "version for the release (required)")
	draftCmd.Flags().StringVar(&draftOpts.tag, "tag", "", "tag for the release")
	draftCmd.Flags().StringVar(&draftOpts.name, "name", "", "human-readable name for the release")
	draftCmd.Flags().StringVar(&draftOpts.description, "description", "", "description for the release (e.g. release notes)")
	draftCmd.Flags().StringVar(&draftOpts.channel, "channel", "stable", "channel for the release, one of: stable, rc, beta, alpha, dev")
	draftCmd.Flags().BoolVar(&draftOpts.noAutoUpgrade, "no-auto-upgrade", false, "disable automatic upgrade checks [$KEYGEN_NO_AUTO_UPGRADE=1]")

	// TODO(ezekg) Accept entitlement codes and entitlement IDs?
	draftCmd.Flags().StringSliceVar(&draftOpts.entitlements, "entitlements", []string{}, "comma seperated list of entitlement constraints (e.g. --entitlements <id>,<id>,...)")

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

	if v, ok := os.LookupEnv("KEYGEN_PRODUCT_TOKEN"); ok {
		if keygenext.Token == "" {
			keygenext.Token = v
		}
	}

	if _, ok := os.LookupEnv("KEYGEN_NO_AUTO_UPGRADE"); ok {
		draftOpts.noAutoUpgrade = true
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
	if !draftOpts.noAutoUpgrade {
		err := upgradeRun(nil, nil)
		if err != nil {
			return err
		}
	}

	italic := color.New(color.Italic).SprintFunc()
	channel := draftOpts.channel

	var constraints keygenext.Constraints
	if e := draftOpts.entitlements; len(e) != 0 {
		constraints = constraints.From(e)
	}

	var tag *string
	if t := draftOpts.tag; t != "" {
		tag = &t
	}

	var name *string
	if n := draftOpts.name; n != "" {
		name = &n
	}

	var desc *string
	if d := draftOpts.description; d != "" {
		desc = &d
	}

	version, err := semver.NewVersion(draftOpts.version)
	if err != nil {
		return fmt.Errorf(`version "%s" is not acceptable (%s)`, draftOpts.version, italic(strings.ToLower(err.Error())))
	}

	release := &keygenext.Release{
		Name:        name,
		Description: desc,
		Version:     version.String(),
		Tag:         tag,
		Channel:     channel,
		ProductID:   keygenext.Product,
		Constraints: constraints,
	}

	if err := release.Create(); err != nil {
		if e, ok := err.(*keygenext.Error); ok {
			var code string
			if e.Code != "" {
				code = italic("(" + e.Code + ")")
			}

			return fmt.Errorf("%s: %s %s", e.Title, e.Detail, code)
		}

		return err
	}

	fmt.Println("drafted release " + italic(release.ID))

	return nil
}
