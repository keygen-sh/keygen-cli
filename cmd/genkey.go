package cmd

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
	"github.com/spf13/cobra"
)

var (
	genkeyOpts = &CommandOptions{}
	genkeyCmd  = &cobra.Command{
		Use:   "genkey",
		Short: "generate an ed25519 key pair for code signing",
		Args:  cobra.NoArgs,
		RunE:  genkeyRun,

		// Encountering an error should not display usage
		SilenceUsage: true,
	}
)

func init() {
	genkeyCmd.Flags().StringVar(&genkeyOpts.signingKeyPath, "out", "keygen.key", "output the private publishing key to specified file")
	genkeyCmd.Flags().StringVar(&genkeyOpts.verifyKeyPath, "pubout", "keygen.pub", "output the public upgrade key to specified file")
	genkeyCmd.Flags().BoolVar(&genkeyOpts.noAutoUpgrade, "no-auto-upgrade", false, "disable automatic upgrade checks [$KEYGEN_NO_AUTO_UPGRADE=1]")

	if _, ok := os.LookupEnv("KEYGEN_NO_AUTO_UPGRADE"); ok {
		genkeyOpts.noAutoUpgrade = true
	}

	rootCmd.AddCommand(genkeyCmd)
}

func genkeyRun(cmd *cobra.Command, args []string) error {
	if !genkeyOpts.noAutoUpgrade {
		err := upgradeRun(nil, nil)
		if err != nil {
			return err
		}
	}

	signingKeyPath, err := homedir.Expand(genkeyOpts.signingKeyPath)
	if err != nil {
		return fmt.Errorf(`path "%s" is not expandable (%s)`, genkeyOpts.signingKeyPath, err)
	}

	verifyKeyPath, err := homedir.Expand(genkeyOpts.verifyKeyPath)
	if err != nil {
		return fmt.Errorf(`path "%s" is not expandable (%s)`, genkeyOpts.verifyKeyPath, err)
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

	yellow := color.New(color.FgYellow).SprintFunc()
	italic := color.New(color.Italic).SprintFunc()

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
