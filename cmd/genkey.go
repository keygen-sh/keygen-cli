package cmd

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/keygen-sh/keygen-cli/internal/ed25519ph"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	genkeyOpts = &CommandOptions{}
	genkeyCmd  = &cobra.Command{
		Use:   "genkey",
		Short: "generate an ed25519 keypair for code signing",
		RunE:  genkeyRun,

		// Encountering an error should not display usage
		SilenceUsage: true,
	}
)

func init() {
	genkeyCmd.Flags().StringVar(&genkeyOpts.signingKey, "signing-key", "keygen.key", "write private signing key to file")
	genkeyCmd.Flags().StringVar(&genkeyOpts.verifyKey, "verify-key", "keygen.pub", "write public verify key to file")

	rootCmd.AddCommand(genkeyCmd)
}

func genkeyRun(cmd *cobra.Command, args []string) error {
	signingKeyPath, err := homedir.Expand(genkeyOpts.signingKey)
	if err != nil {
		return fmt.Errorf(`path "%s" is not expandable (%s)`, genkeyOpts.signingKey, err)
	}

	verifyKeyPath, err := homedir.Expand(genkeyOpts.verifyKey)
	if err != nil {
		return fmt.Errorf(`path "%s" is not expandable (%s)`, genkeyOpts.verifyKey, err)
	}

	if _, err := os.Stat(signingKeyPath); err == nil {
		return fmt.Errorf(`signing key file "%s" already exists`, signingKeyPath)
	}

	if _, err := os.Stat(verifyKeyPath); err == nil {
		return fmt.Errorf(`verify key file "%s" already exists`, verifyKeyPath)
	}

	verifyKey, signingKey, err := ed25519ph.GenerateKey()
	if err != nil {
		return err
	}

	if err := writesigningKeyFile(signingKeyPath, signingKey); err != nil {
		return err
	}

	if err := writeverifyKeyFile(verifyKeyPath, verifyKey); err != nil {
		return err
	}

	if abs, err := filepath.Abs(signingKeyPath); err != nil {
		fmt.Printf("Signing key: %s\n", signingKeyPath)
	} else {
		fmt.Printf("Signing key: %s\n", abs)
	}

	if abs, err := filepath.Abs(verifyKeyPath); err != nil {
		fmt.Printf("Verify key: %s\n", verifyKeyPath)
	} else {
		fmt.Printf("Verify key: %s\n", abs)
	}

	return nil
}

func writesigningKeyFile(signingKeyPath string, signingKey ed25519ph.SigningKey) error {
	file, err := os.Create(signingKeyPath)
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

func writeverifyKeyFile(verifyKeyPath string, verifyKey ed25519ph.VerifyKey) error {
	file, err := os.Create(verifyKeyPath)
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
