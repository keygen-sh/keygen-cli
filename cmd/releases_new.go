package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/keygen-sh/cli/keygenext"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	releasesNewCmd = &cobra.Command{
		Use:   "new <path>",
		Short: "Publish a new release for a product",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("path to file is required")
			}

			path, err := homedir.Expand(args[0])
			if err != nil {
				return fmt.Errorf(`path "%s" is not expandable (%s)`, args[0], err)
			}

			info, err := os.Stat(path)
			if err != nil {
				reason, ok := err.(*os.PathError)
				if !ok {
					return err
				}

				return fmt.Errorf(`path "%s" is not readable (%s)`, path, reason.Err)
			}

			if info.IsDir() {
				return fmt.Errorf(`path "%s" is a directory (must be a file)`, path)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer file.Close()

			info, err := file.Stat()
			if err != nil {
				return err
			}

			filename := file.Name()
			filesize := info.Size()
			filetype := filepath.Ext(filename)
			if filetype == "" {
				filetype = "bin"
			}

			channel := flags.channel
			platform := flags.platform
			constraints := flags.constraints

			var name *string
			if n := flags.name; n != "" {
				name = &n
			}

			// TODO(ezekg) Transform entitlement codes to entitlement IDs
			// entitlements, err := getEntitlements(constraints)

			version, err := semver.NewVersion(flags.version)
			if err != nil {
				return err
			}

			signature := flags.signature
			if signature == "" {
				// TODO(ezekg) Sign release
			}

			checksum := flags.checksum
			if checksum == "" {
				// TODO(ezekg) Hash release
			}

			release := &keygenext.Release{
				Name:      name,
				Version:   version.String(),
				Filename:  filename,
				Filesize:  int(filesize),
				Filetype:  filetype,
				Platform:  platform,
				Channel:   channel,
				ProductID: keygenext.Product,
				Constraints: keygenext.Constraints{
					{EntitlementID: constraints[0]},
					{EntitlementID: constraints[1]},
					{EntitlementID: constraints[2]},
				},
			}

			fmt.Println(release.Upsert())
			fmt.Println(release)

			return nil
		},
	}
)

func init() {
	f := releasesNewCmd.Flags()

	f.StringVar(&flags.version, "version", "", "version for the release (required)")
	f.StringVar(&flags.name, "name", "", "name for the release")
	f.StringVar(&flags.platform, "platform", "", "platform for the release (required)")
	f.StringVar(&flags.channel, "channel", "stable", "channel for the release, one of: stable, rc, beta, alpha, dev")
	f.StringVar(&flags.signature, "signature", "", "precalculated signature for the release (will be signed using ed25519 by default)")
	f.StringVar(&flags.checksum, "checksum", "", "precalculated checksum for the release (will be hashed using sha-512 by default)")
	f.StringSliceVar(&flags.constraints, "constraints", []string{}, "comma seperated list of entitlement IDs or codes (e.g. --constraints ENTL1,ENTL2,ENTL3)")

	// TODO(ezekg) Add signing key flag
	// TODO(ezekg) Add metadata flag

	releasesNewCmd.MarkPersistentFlagRequired("version")
	releasesNewCmd.MarkPersistentFlagRequired("platform")

	releasesCmd.AddCommand(releasesNewCmd)
}
