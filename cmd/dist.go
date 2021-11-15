package cmd

import (
	"crypto"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/keygen-sh/keygen-cli/internal/keygenext"
	"github.com/mitchellh/go-homedir"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
	"github.com/spf13/cobra"
)

var (
	distOpts = &CommandOptions{}
	distCmd  = &cobra.Command{
		Use:   "dist <path>",
		Short: "publish a new release for a product",
		Args:  distArgs,
		RunE:  distRun,

		// Encountering an error should not display usage
		SilenceUsage: true,
	}
)

func init() {
	distCmd.Flags().StringVar(&keygenext.Account, "account", "", "your keygen.sh account identifier [$KEYGEN_ACCOUNT_ID] (required)")
	distCmd.Flags().StringVar(&keygenext.Product, "product", "", "your keygen.sh product identifier [$KEYGEN_PRODUCT_ID] (required)")
	distCmd.Flags().StringVar(&keygenext.Token, "token", "", "your keygen.sh product token [$KEYGEN_PRODUCT_TOKEN] (required)")
	distCmd.Flags().StringVar(&distOpts.filename, "filename", "", "filename for the release (default is filename from <path>)")
	distCmd.Flags().StringVar(&distOpts.version, "version", "", "version for the release (required)")
	distCmd.Flags().StringVar(&distOpts.name, "name", "", "human-readable name for the release")
	distCmd.Flags().StringVar(&distOpts.description, "description", "", "description for the release (e.g. release notes)")
	distCmd.Flags().StringVar(&distOpts.platform, "platform", "", "platform for the release (required)")
	distCmd.Flags().StringVar(&distOpts.channel, "channel", "stable", "channel for the release, one of: stable, rc, beta, alpha, dev")
	distCmd.Flags().StringVar(&distOpts.signature, "signature", "", "pre-calculated signature for the release (defaults using ed25519ph)")
	distCmd.Flags().StringVar(&distOpts.checksum, "checksum", "", "pre-calculated checksum for the release (defaults using sha-512)")
	distCmd.Flags().StringVar(&distOpts.signingAlgorithm, "signing-algorithm", "ed25519ph", "the signing algorithm to use, one of: ed25519ph, ed25519")
	distCmd.Flags().StringVar(&distOpts.signingKey, "signing-key", "", "path to ed25519 private key for signing the release [$KEYGEN_SIGNING_KEY]")

	// TODO(ezekg) Accept entitlement codes and entitlement IDs?
	distCmd.Flags().StringSliceVar(&distOpts.entitlements, "entitlements", []string{}, "comma seperated list of entitlement constraints (e.g. --entitlements <id>,<id>,...)")

	// TODO(ezekg) Prompt multi-line description input from stdin if "--"?
	// TODO(ezekg) Add metadata flag

	if v := os.Getenv("KEYGEN_ACCOUNT_ID"); v != "" {
		if keygenext.Account == "" {
			keygenext.Account = v
		}
	}

	if v := os.Getenv("KEYGEN_PRODUCT_ID"); v != "" {
		if keygenext.Product == "" {
			keygenext.Product = v
		}
	}

	if v := os.Getenv("KEYGEN_PRODUCT_TOKEN"); v != "" {
		if keygenext.Token == "" {
			keygenext.Token = v
		}
	}

	if v := os.Getenv("KEYGEN_SIGNING_KEY"); v != "" {
		if distOpts.signingKey == "" {
			distOpts.signingKey = v
		}
	}

	if keygenext.Account == "" {
		distCmd.MarkFlagRequired("account")
	}

	if keygenext.Product == "" {
		distCmd.MarkFlagRequired("product")
	}

	if keygenext.Token == "" {
		distCmd.MarkFlagRequired("token")
	}

	distCmd.MarkFlagRequired("version")
	distCmd.MarkFlagRequired("platform")

	rootCmd.AddCommand(distCmd)
}

func distArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("path to file is required")
	}

	return nil
}

func distRun(cmd *cobra.Command, args []string) error {
	// TODO(ezekg) Add no-auto-upgrade flag
	err := upgradeRun(nil, nil)
	if err != nil {
		return err
	}

	path, err := homedir.Expand(args[0])
	if err != nil {
		return fmt.Errorf(`path "%s" is not expandable (%s)`, args[0], err)
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf(`path "%s" is not readable (%s)`, path, err.(*os.PathError).Err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf(`path "%s" is not readable (%s)`, path, err.(*os.PathError).Err)
	}

	if info.IsDir() {
		return fmt.Errorf(`path "%s" is a directory (must be a file)`, path)
	}

	filename := file.Name()
	filesize := info.Size()
	filetype := filepath.Ext(filename)
	if filetype == "" {
		filetype = "bin"
	}

	// Allow filename to be overridden
	if n := distOpts.filename; n != "" {
		filename = n
	}

	channel := distOpts.channel
	platform := distOpts.platform

	constraints := keygenext.Constraints{}
	if e := distOpts.entitlements; len(e) != 0 {
		constraints = constraints.From(e)
	}

	var name *string
	if n := distOpts.name; n != "" {
		name = &n
	}

	var desc *string
	if d := distOpts.description; d != "" {
		desc = &d
	}

	version, err := semver.NewVersion(distOpts.version)
	if err != nil {
		return fmt.Errorf(`version "%s" is not acceptable (%s)`, distOpts.version, strings.ToLower(err.Error()))
	}

	checksum := distOpts.checksum
	if checksum == "" {
		checksum, err = calculateChecksum(file)
		if err != nil {
			return err
		}
	}

	signature := distOpts.signature
	if pk := distOpts.signingKey; pk != "" && signature == "" {
		path, err := homedir.Expand(pk)
		if err != nil {
			return fmt.Errorf(`signing-key "%s" is not expandable (%s)`, pk, err)
		}

		signature, err = calculateSignature(path, file)
		if err != nil {
			return err
		}
	}

	release := &keygenext.Release{
		Name:        name,
		Description: desc,
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
	if err := release.Upsert(); err != nil {
		return err
	}

	if err := release.Upload(file); err != nil {
		return err
	}

	fmt.Println(`successfully published release "` + release.ID + `"`)

	return nil
}

func calculateChecksum(file *os.File) (string, error) {
	defer file.Seek(0, io.SeekStart) // reset reader

	h := sha512.New()

	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}

	digest := h.Sum(nil)

	return base64.RawStdEncoding.EncodeToString(digest), nil
}

func calculateSignature(signingKeyPath string, file *os.File) (string, error) {
	defer file.Seek(0, io.SeekStart) // reset reader

	var signingKey ed25519.PrivateKey
	encKey, err := os.ReadFile(signingKeyPath)
	if err != nil {
		return "", err
	}

	signingKey, err = hex.DecodeString(string(encKey))
	if err != nil {
		return "", err
	}

	if l := len(signingKey); l != ed25519.PrivateKeySize {
		return "", fmt.Errorf("bad signing key length (got %d expected 64)", l)
	}

	var sig []byte
	switch distOpts.signingAlgorithm {
	case "ed25519ph":
		// We're using Ed25519ph which expects a pre-hashed message using SHA-512
		h := sha512.New()

		if _, err := io.Copy(h, file); err != nil {
			return "", err
		}

		opts := &ed25519.Options{Hash: crypto.SHA512, Context: keygenext.Product}
		digest := h.Sum(nil)

		sig, err = signingKey.Sign(nil, digest, opts)
		if err != nil {
			return "", err
		}
	case "ed25519":
		fmt.Println("warning: using ed25519 to sign large files is not recommended (use ed25519ph instead)")

		b, err := ioutil.ReadAll(file)
		if err != nil {
			return "", err
		}

		// TODO(ezekg) Validate key size to guard against Sign panicing
		sig, err = signingKey.Sign(nil, b, &ed25519.Options{})
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf(`signing algorithm "%s" is not supported`, distOpts.signingAlgorithm)
	}

	return base64.RawStdEncoding.EncodeToString(sig), nil
}
