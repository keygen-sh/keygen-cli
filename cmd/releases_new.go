package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/keygen-sh/cli/keygenext"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	releasesNewCmd = &cobra.Command{
		Use:   "new <path>",
		Short: "Publish a new release for a product",
		Args:  releasesNewArgs,
		RunE:  releasesNewRun,
	}
)

func init() {
	releasesNewCmd.Flags().StringVar(&flags.version, "version", "", "version for the release (required)")
	releasesNewCmd.Flags().StringVar(&flags.name, "name", "", "name for the release")
	releasesNewCmd.Flags().StringVar(&flags.platform, "platform", "", "platform for the release (required)")
	releasesNewCmd.Flags().StringVar(&flags.channel, "channel", "stable", "channel for the release, one of: stable, rc, beta, alpha, dev")
	releasesNewCmd.Flags().StringVar(&flags.signature, "signature", "", "precalculated signature for the release (will be signed using ed25519 by default)")
	releasesNewCmd.Flags().StringVar(&flags.checksum, "checksum", "", "precalculated checksum for the release (will be hashed using sha-512 by default)")
	releasesNewCmd.Flags().StringVar(&flags.signingKey, "signing-key", "", "path to the ed25519 private key for signing releases")

	// TODO(ezekg) Accept entitlement codes and entitlement IDs?
	releasesNewCmd.Flags().StringSliceVar(&flags.constraints, "constraints", []string{}, "comma seperated list of entitlement IDs (e.g. --constraints <id>,<id>,...)")

	// TODO(ezekg) Add metadata flag

	releasesNewCmd.MarkFlagRequired("version")
	releasesNewCmd.MarkFlagRequired("platform")

	releasesCmd.AddCommand(releasesNewCmd)
}

func releasesNewArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("path to file is required")
	}

	path, err := homedir.Expand(args[0])
	if err != nil {
		return fmt.Errorf(`path "%s" is not expandable (%s)`, args[0], err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf(`path "%s" is not readable (%s)`, path, err.(*os.PathError).Err)
	}

	if info.IsDir() {
		return fmt.Errorf(`path "%s" is a directory (must be a file)`, path)
	}

	return nil
}

func releasesNewRun(cmd *cobra.Command, args []string) error {
	path := args[0]
	file, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf(`path "%s" is not readable (%s)`, path, err.(*os.PathError).Err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf(`path "%s" is not readable (%s)`, path, err.(*os.PathError).Err)
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

	version, err := semver.NewVersion(flags.version)
	if err != nil {
		return fmt.Errorf(`version "%s" is not acceptable (%s)`, flags.version, strings.ToLower(err.Error()))
	}

	var signature *string
	if sig := flags.signature; sig != "" {
		signature = &sig
	}

	if key := flags.signingKey; signature == nil && key != "" {
		// TODO(ezekg) Sign release
	}

	var checksum *string
	if c := flags.checksum; c == "" {
		// TODO(ezekg) Hash release
	} else {
		checksum = &c
	}

	release := &keygenext.Release{
		Name:      name,
		Version:   version.String(),
		Filename:  filename,
		Filesize:  filesize,
		Filetype:  filetype,
		Platform:  platform,
		Signature: signature,
		Checksum:  checksum,
		Channel:   channel,
		ProductID: keygenext.Product,
		Constraints: keygenext.Constraints{
			{EntitlementID: constraints[0]},
			{EntitlementID: constraints[1]},
			{EntitlementID: constraints[2]},
		},
	}

	err = release.Upsert()
	if err != nil {
		return err
	}

	err = release.Upload(file)
	if err != nil {
		return err
	}

	return nil
}
