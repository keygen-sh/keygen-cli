package cmd

import (
	"bufio"
	"crypto"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/keygen-sh/keygen-cli/internal/keygenext"
	"github.com/mattn/go-isatty"
	"github.com/mitchellh/go-homedir"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

var (
	uploadOpts = &UploadCommandOptions{}
	uploadCmd  = &cobra.Command{
		Use:   "upload <path>",
		Short: "upload a new artifact for a release",
		Example: `  keygen upload ./build/keygen_darwin_amd64 \
      --signing-key ~/.keys/keygen.key \
      --account '1fddcec8-8dd3-4d8d-9b16-215cac0f9b52' \
      --product '2313b7e7-1ea6-4a01-901e-2931de6bb1e2' \
      --token 'prod-xxx' \
      --release '1.0.0' \
      --platform 'darwin' \
      --arch 'amd64' \
      --metadata '{"key": "value"}'

Docs:
  https://keygen.sh/docs/cli/`,
		Args: uploadArgs,
		RunE: uploadRun,

		// Encountering an error should not display usage
		SilenceUsage: true,
	}
)

type UploadCommandOptions struct {
	Filename          string
	Filetype          string
	Platform          string
	Arch              string
	Release           string
	Package           string
	Checksum          string
	ChecksumAlgorithm string
	ChecksumEncoding  string
	Signature         string
	SigningAlgorithm  string
	SignatureEncoding string
	SigningKeyPath    string
	SigningKey        string
	NoAutoUpgrade     bool
	Metadata          string
}

func init() {
	uploadCmd.Flags().StringVar(&keygenext.Account, "account", "", "your keygen.sh account identifier [$KEYGEN_ACCOUNT_ID=<id>] (required)")
	uploadCmd.Flags().StringVar(&keygenext.Product, "product", "", "your keygen.sh product identifier [$KEYGEN_PRODUCT_ID=<id>] (required)")
	uploadCmd.Flags().StringVar(&keygenext.Token, "token", "", "your keygen.sh product or environment token [$KEYGEN_TOKEN] (required)")
	uploadCmd.Flags().StringVar(&keygenext.Environment, "environment", "", "your keygen.sh environment identifier [$KEYGEN_ENVIRONMENT=<id>]")
	uploadCmd.Flags().StringVar(&keygenext.APIURL, "host", "", "the host of the keygen server [$KEYGEN_HOST=<host>]")
	uploadCmd.Flags().StringVar(&uploadOpts.Release, "release", "", "the release identifier (required)")
	uploadCmd.Flags().StringVar(&uploadOpts.Package, "package", "", "package identifier for the artifact")
	uploadCmd.Flags().StringVar(&uploadOpts.Filename, "filename", "", "filename for the artifact (defaults to basename of <path>)")
	uploadCmd.Flags().StringVar(&uploadOpts.Filetype, "filetype", "<auto>", "filetype for the artifact (defaults to extname of <path>)")
	uploadCmd.Flags().StringVar(&uploadOpts.Platform, "platform", "", "platform for the artifact")
	uploadCmd.Flags().StringVar(&uploadOpts.Arch, "arch", "", "arch for the artifact")
	uploadCmd.Flags().StringVar(&uploadOpts.Checksum, "checksum", "", "pre-calculated checksum for the artifact (defaults using sha-512)")
	uploadCmd.Flags().StringVar(&uploadOpts.ChecksumAlgorithm, "checksum-algorithm", "sha-512", "the checksum algorithm to use, one of: sha-512, sha-256, sha-1")
	uploadCmd.Flags().StringVar(&uploadOpts.ChecksumEncoding, "checksum-encoding", "base64raw", "the checksum encoding to use, one of: base64, base64raw, base64url, hex")
	uploadCmd.Flags().StringVar(&uploadOpts.Signature, "signature", "", "pre-calculated signature for the artifact (defaults using ed25519ph)")
	uploadCmd.Flags().StringVar(&uploadOpts.SigningAlgorithm, "signing-algorithm", "ed25519ph", "the signing algorithm to use, one of: ed25519ph, ed25519")
	uploadCmd.Flags().StringVar(&uploadOpts.SignatureEncoding, "signature-encoding", "base64raw", "the signature encoding to use, one of: base64, base64raw, base64url, hex")
	uploadCmd.Flags().StringVar(&uploadOpts.SigningKeyPath, "signing-key", "", "path to ed25519 private key for signing the artifact [$KEYGEN_SIGNING_KEY_PATH=<path>, $KEYGEN_SIGNING_KEY=<key>]")
	uploadCmd.Flags().BoolVar(&uploadOpts.NoAutoUpgrade, "no-auto-upgrade", false, "disable automatic upgrade checks [$KEYGEN_NO_AUTO_UPGRADE=1]")
	uploadCmd.Flags().StringVar(&uploadOpts.Metadata, "metadata", "", "JSON string of metadata key-value pairs")

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

	if v, ok := os.LookupEnv("KEYGEN_SIGNING_KEY_PATH"); ok {
		if uploadOpts.SigningKeyPath == "" {
			uploadOpts.SigningKeyPath = v
		}
	}

	if v, ok := os.LookupEnv("KEYGEN_SIGNING_KEY"); ok {
		if uploadOpts.SigningKey == "" {
			uploadOpts.SigningKey = v
		}
	}

	if _, ok := os.LookupEnv("KEYGEN_NO_AUTO_UPGRADE"); ok {
		uploadOpts.NoAutoUpgrade = true
	}

	if keygenext.Account == "" {
		uploadCmd.MarkFlagRequired("account")
	}

	if keygenext.Product == "" {
		uploadCmd.MarkFlagRequired("product")
	}

	if keygenext.Token == "" {
		uploadCmd.MarkFlagRequired("token")
	}

	uploadCmd.MarkFlagRequired("release")

	rootCmd.AddCommand(uploadCmd)
}

func uploadArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("path is required")
	}

	return nil
}

func uploadRun(cmd *cobra.Command, args []string) error {
	if !uploadOpts.NoAutoUpgrade {
		err := upgradeRun(nil, nil)
		if err != nil {
			return err
		}
	}

	path, err := homedir.Expand(args[0])
	if err != nil {
		return fmt.Errorf(`path "%s" is not expandable (%s)`, args[0], italic(err))
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf(`path "%s" is not readable (%s)`, path, italic(err.(*os.PathError).Err))
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf(`path "%s" is not readable (%s)`, path, italic(err.(*os.PathError).Err))
	}

	if info.IsDir() {
		return fmt.Errorf(`path "%s" is a directory (must be a file)`, path)
	}

	platform := uploadOpts.Platform
	arch := uploadOpts.Arch
	filename := filepath.Base(info.Name())
	filesize := info.Size()

	// Allow filename to be overridden
	if n := uploadOpts.Filename; n != "" {
		filename = n
	}

	// Allow filetype to be overridden
	var filetype string

	if uploadOpts.Filetype == "<auto>" {
		filetype = filepath.Ext(filename)
		if _, e := strconv.Atoi(filetype); e == nil {
			filetype = ""
		}
	} else {
		filetype = uploadOpts.Filetype
	}

	checksum := uploadOpts.Checksum
	if checksum == "" {
		checksum, err = calculateChecksum(file)
		if err != nil {
			return err
		}
	}

	signature := uploadOpts.Signature
	if signature == "" && (uploadOpts.SigningKeyPath != "" || uploadOpts.SigningKey != "") {
		var key string

		switch {
		case uploadOpts.SigningKeyPath != "":
			path, err := homedir.Expand(uploadOpts.SigningKeyPath)
			if err != nil {
				return fmt.Errorf(`signing-key path is not expandable (%s)`, err)
			}

			b, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf(`signing-key path is not readable (%s)`, err)
			}

			key = string(b)
		case uploadOpts.SigningKey != "":
			key = uploadOpts.SigningKey
		}

		signature, err = calculateSignature(key, file)
		if err != nil {
			return err
		}
	}

	var metadata map[string]interface{}
	if m := uploadOpts.Metadata; m != "" {
		if err := json.Unmarshal([]byte(m), &metadata); err != nil {
			return fmt.Errorf("invalid metadata JSON: %v", err)
		}
	}

	release := &keygenext.Release{
		ID:        uploadOpts.Release,
		PackageID: &uploadOpts.Package,
		ProductID: keygenext.Product,
	}

	// get actual release id w/ filters e.g. package
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

	artifact := &keygenext.Artifact{
		Filename:  filename,
		Filesize:  filesize,
		Filetype:  filetype,
		Platform:  platform,
		Arch:      arch,
		Signature: signature,
		Checksum:  checksum,
		ReleaseID: &release.ID,
		Metadata:  metadata,
	}
	if err := artifact.Create(); err != nil {
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

	// Create a buffered reader to limit memory footprint
	var reader io.Reader = bufio.NewReaderSize(file, 1024*1024*50 /* 50 mb */)
	var progress *mpb.Progress

	// Create a progress bar for file upload if TTY
	if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		progress = mpb.New(mpb.WithWidth(60), mpb.WithRefreshRate(180*time.Millisecond))
		bar := progress.Add(
			artifact.Filesize,
			mpb.NewBarFiller(mpb.BarStyle().Rbound("|")),
			mpb.BarRemoveOnComplete(),
			mpb.PrependDecorators(
				decor.CountersKibiByte("% .2f / % .2f"),
			),
			mpb.AppendDecorators(
				decor.EwmaETA(decor.ET_STYLE_GO, 90),
				decor.Name(" ] "),
				decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
			),
		)

		// Create proxy reader for the progress bar
		reader = bar.ProxyReader(reader)
		closer, ok := reader.(io.ReadCloser)
		if ok {
			defer closer.Close()
		}
	}

	if err := artifact.Upload(reader); err != nil {
		return err
	}

	if progress != nil {
		progress.Wait()
	}

	fmt.Println(green("uploaded:") + " artifact " + italic(artifact.ID))

	return nil
}

func calculateChecksum(file *os.File) (string, error) {
	defer file.Seek(0, io.SeekStart) // reset reader
	var h hash.Hash

	switch uploadOpts.ChecksumAlgorithm {
	case "sha-512":
		h = sha512.New()
	case "sha-256":
		h = sha256.New()
	case "sha-1":
		h = sha1.New()
	default:
		return "", fmt.Errorf(`checksum algorithm "%s" is not supported`, uploadOpts.ChecksumAlgorithm)
	}

	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}

	digest := h.Sum(nil)

	switch uploadOpts.ChecksumEncoding {
	case "base64":
		return base64.StdEncoding.EncodeToString(digest), nil
	case "base64raw":
		return base64.RawStdEncoding.EncodeToString(digest), nil
	case "base64url":
		return base64.URLEncoding.EncodeToString(digest), nil
	case "hex":
		return hex.EncodeToString(digest), nil
	default:
		return "", fmt.Errorf(`checksum encoding "%s" is not supported`, uploadOpts.ChecksumEncoding)
	}
}

func calculateSignature(encSigningKey string, file *os.File) (string, error) {
	defer file.Seek(0, io.SeekStart) // reset reader

	decSigningKey, err := hex.DecodeString(encSigningKey)
	if err != nil {
		return "", fmt.Errorf("bad signing key (%s)", err)
	}

	if l := len(decSigningKey); l != ed25519.PrivateKeySize {
		return "", fmt.Errorf("bad signing key length (got %d expected %d)", l, ed25519.PrivateKeySize)
	}

	signingKey := ed25519.PrivateKey(decSigningKey)
	var sig []byte

	switch uploadOpts.SigningAlgorithm {
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
		fmt.Println(yellow("warning:") + " using ed25519 to sign large files is not recommended (use ed25519ph instead)")

		b, err := ioutil.ReadAll(file)
		if err != nil {
			return "", err
		}

		sig, err = signingKey.Sign(nil, b, &ed25519.Options{})
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf(`signing algorithm "%s" is not supported`, uploadOpts.SigningAlgorithm)
	}

	switch uploadOpts.SignatureEncoding {
	case "base64":
		return base64.StdEncoding.EncodeToString(sig), nil
	case "base64raw":
		return base64.RawStdEncoding.EncodeToString(sig), nil
	case "base64url":
		return base64.URLEncoding.EncodeToString(sig), nil
	case "hex":
		return hex.EncodeToString(sig), nil
	default:
		return "", fmt.Errorf(`signature encoding "%s" is not supported`, uploadOpts.SignatureEncoding)
	}
}
