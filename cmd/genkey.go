package cmd

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
	"github.com/spf13/cobra"
)

var (
	genkeyOpts = &GenKeyCommandOptions{}
	genkeyCmd  = &cobra.Command{
		Use:   "genkey",
		Short: "generate an ed25519 key pair for code signing",
		Args:  cobra.NoArgs,
		RunE:  genkeyRun,

		// Encountering an error should not display usage
		SilenceUsage: true,
	}
)

type GenKeyCommandOptions struct {
	SigningKeyPath string
	VerifyKeyPath  string
	NoAutoUpgrade  bool
}

func init() {
	genkeyCmd.Flags().StringVar(&genkeyOpts.SigningKeyPath, "out", "keygen.key", "output the private publishing key to specified file")
	genkeyCmd.Flags().StringVar(&genkeyOpts.VerifyKeyPath, "pubout", "keygen.pub", "output the public upgrade key to specified file")
	genkeyCmd.Flags().BoolVar(&genkeyOpts.NoAutoUpgrade, "no-auto-upgrade", false, "disable automatic upgrade checks [$KEYGEN_NO_AUTO_UPGRADE=1]")

	if _, ok := os.LookupEnv("KEYGEN_NO_AUTO_UPGRADE"); ok {
		genkeyOpts.NoAutoUpgrade = true
	}

	rootCmd.AddCommand(genkeyCmd)
}

func genkeyRun(cmd *cobra.Command, args []string) error {
	if !genkeyOpts.NoAutoUpgrade {
		err := upgradeRun(nil, nil)
		if err != nil {
			return err
		}
	}

	signingKeyPath, err := homedir.Expand(genkeyOpts.SigningKeyPath)
	if err != nil {
		return fmt.Errorf(`path "%s" is not expandable (%s)`, genkeyOpts.SigningKeyPath, err)
	}

	verifyKeyPath, err := homedir.Expand(genkeyOpts.VerifyKeyPath)
	if err != nil {
		return fmt.Errorf(`path "%s" is not expandable (%s)`, genkeyOpts.VerifyKeyPath, err)
	}

	if _, err := os.Stat(signingKeyPath); err == nil {
		return fmt.Errorf(`private key file "%s" already exists`, signingKeyPath)
	}

	if _, err := os.Stat(verifyKeyPath); err == nil {
		return fmt.Errorf(`public key file "%s" already exists`, verifyKeyPath)
	}

	verifyKey, signingKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}

	if err := writeSigningKeyFile(signingKeyPath, signingKey); err != nil {
		return err
	}

	if err := writeVerifyKeyFile(verifyKeyPath, verifyKey); err != nil {
		return err
	}

	if abs, err := filepath.Abs(signingKeyPath); err == nil {
		signingKeyPath = abs
	}

	if abs, err := filepath.Abs(verifyKeyPath); err == nil {
		verifyKeyPath = abs
	}

	fmt.Printf(`private key: %s
public key: %s
`,
		signingKeyPath,
		verifyKeyPath,
	)

	fmt.Fprintf(os.Stderr, yellow("warning:")+" never share your private key -- "+italic("it's a secret!")+"\n")

	return nil
}

func writeSigningKeyFile(signingKeyPath string, signingKey ed25519.PrivateKey) error {
	file, err := os.OpenFile(signingKeyPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := hex.EncodeToString(signingKey)
	_, err = file.WriteString(enc)
	if err != nil {
		return err
	}

	return nil
}

func writeVerifyKeyFile(verifyKeyPath string, verifyKey ed25519.PublicKey) error {
	file, err := os.OpenFile(verifyKeyPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := hex.EncodeToString(verifyKey)
	_, err = file.WriteString(enc)
	if err != nil {
		return err
	}

	return nil
}
