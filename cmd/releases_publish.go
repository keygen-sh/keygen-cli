package cmd

import (
	"bufio"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/keygen-sh/keygen-cli/internal/ed25519ph"
	"github.com/keygen-sh/keygen-cli/internal/keygenext"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	releasesPublishCmd = &cobra.Command{
		Use:   "publish <path>",
		Short: "publish a new release for a product",
		Args:  releasesPublishArgs,
		RunE:  releasesPublishRun,

		// Encountering an error should not display usage
		SilenceUsage: true,
	}
)

func init() {
	releasesPublishCmd.Flags().StringVar(&flags.filename, "filename", "", "filename for the release (default is filename from <path>)")
	releasesPublishCmd.Flags().StringVar(&flags.version, "version", "", "version for the release (required)")
	releasesPublishCmd.Flags().StringVar(&flags.name, "name", "", "human-readable name for the release")
	releasesPublishCmd.Flags().StringVar(&flags.platform, "platform", "", "platform for the release (required)")
	releasesPublishCmd.Flags().StringVar(&flags.channel, "channel", "stable", "channel for the release, one of: stable, rc, beta, alpha, dev")
	releasesPublishCmd.Flags().StringVar(&flags.signingKey, "signing-key", "", "hex-encoded ed25519 private key for signing releases")

	// TODO(ezekg) Accept entitlement codes and entitlement IDs?
	releasesPublishCmd.Flags().StringSliceVar(&flags.constraints, "constraints", []string{}, "comma seperated list of entitlement identifiers (e.g. --constraints <id>,<id>,...)")

	// TODO(ezekg) Add metadata flag

	releasesPublishCmd.MarkFlagRequired("version")
	releasesPublishCmd.MarkFlagRequired("platform")

	releasesCmd.AddCommand(releasesPublishCmd)
}

func releasesPublishArgs(cmd *cobra.Command, args []string) error {
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

func releasesPublishRun(cmd *cobra.Command, args []string) error {
	// s := spinner.New(spinner.CharSets[13], 100*time.Millisecond)
	// s.HideCursor = true
	// s.Start()

	path := args[0]
	file, err := os.Open(path)
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

	// Allow filename to be overridden
	if n := flags.filename; n != "" {
		filename = n
	}

	channel := flags.channel
	platform := flags.platform

	constraints := keygenext.Constraints{}
	if c := flags.constraints; len(c) != 0 {
		constraints = constraints.From(c)
	}

	var name *string
	if n := flags.name; n != "" {
		name = &n
	}

	version, err := semver.NewVersion(flags.version)
	if err != nil {
		return fmt.Errorf(`version "%s" is not acceptable (%s)`, flags.version, strings.ToLower(err.Error()))
	}

	// s.Suffix = " generating checksum for release..."

	// Reopen file so that we don't screw with our first reader
	checksum, err := calculateChecksum(file)
	if err != nil {
		return err
	}

	var signature *string
	if sk := flags.signingKey; sk != "" {
		// s.Suffix = " generating signature for release..."

		sig, err := calculateSignature(sk, file)
		if err != nil {
			return err
		}

		signature = &sig
	}

	release := &keygenext.Release{
		Name:        name,
		Version:     version.String(),
		Filename:    filename,
		Filesize:    filesize,
		Filetype:    filetype,
		Platform:    platform,
		Signature:   signature,
		Checksum:    checksum,
		Channel:     channel,
		ProductID:   keygenext.Product,
		Constraints: constraints,
	}

	// TODO(ezekg) Should we do a Create() unless a --upsert flag is given?
	// s.Suffix = " creating release object..."

	if err := release.Upsert(); err != nil {
		return err
	}

	// s.Suffix = " uploading release artifact..."

	if err := release.Upload(file); err != nil {
		return err
	}

	// s.Suffix = ""
	// s.FinalMSG = "successfully published release \"" + release.ID + "\"\n"
	// s.Stop()

	return nil
}

func calculateChecksum(file *os.File) (string, error) {
	defer file.Seek(0, io.SeekStart) // reset reader

	reader := bufio.NewReader(file)
	h := sha512.New()

	_, err := io.Copy(h, reader)
	if err != nil {
		return "", err
	}

	shasum := h.Sum(nil)

	return hex.EncodeToString(shasum), nil
}

func calculateSignature(signingKey string, file *os.File) (string, error) {
	defer file.Seek(0, io.SeekStart) // reset reader

	priv, err := hex.DecodeString(signingKey)
	if err != nil {
		return "", err
	}

	// We're using Ed25519ph which expects a pre-hashed message using SHA-512
	h := sha512.New()

	if _, err = io.Copy(h, file); err != nil {
		return "", err
	}

	digest := h.Sum(nil)

	// TODO(ezekg) Validate key size to guard against Sign panicing
	sig, err := ed25519ph.Sign(priv, digest)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sig), nil
}
