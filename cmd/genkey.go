package cmd

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/keygen-sh/keygen-cli/internal/ed25519ph"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	genkeyCmd = &cobra.Command{
		Use:   "genkey",
		Short: "generate an ed25519 keypair for code signing",
		RunE:  genkeyRun,

		// Encountering an error should not display usage
		SilenceUsage: true,
	}
)

func init() {
	genkeyCmd.Flags().StringVar(&options.privateKey, "private-key", "keygen.key", "write private key to file (e.g. --private-key ~/keygen.key)")
	genkeyCmd.Flags().StringVar(&options.publicKey, "public-key", "keygen.pub", "write public key to file (e.g. --public-key ~/keygen.pub)")

	rootCmd.AddCommand(genkeyCmd)
}

func genkeyRun(cmd *cobra.Command, args []string) error {
	privateKeyPath, err := homedir.Expand(options.privateKey)
	if err != nil {
		return fmt.Errorf(`path "%s" is not expandable (%s)`, options.privateKey, err)
	}

	publicKeyPath, err := homedir.Expand(options.publicKey)
	if err != nil {
		return fmt.Errorf(`path "%s" is not expandable (%s)`, options.publicKey, err)
	}

	if _, err := os.Stat(privateKeyPath); err == nil {
		return fmt.Errorf(`private key file "%s" already exists`, privateKeyPath)
	}

	if _, err := os.Stat(publicKeyPath); err == nil {
		return fmt.Errorf(`public key file "%s" already exists`, publicKeyPath)
	}

	publicKey, privateKey, err := ed25519ph.GenerateKey()
	if err != nil {
		return err
	}

	if err := writePrivateKeyFile(privateKeyPath, privateKey); err != nil {
		return err
	}

	if err := writePublicKeyFile(publicKeyPath, publicKey); err != nil {
		return err
	}

	return nil
}

func writePrivateKeyFile(path string, privateKey ed25519ph.PrivateKey) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := hex.EncodeToString(privateKey)
	_, err = file.WriteString(enc)
	if err != nil {
		return err
	}

	return nil
}

func writePublicKeyFile(path string, publicKey ed25519ph.PublicKey) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := hex.EncodeToString(publicKey)
	_, err = file.WriteString(enc)
	if err != nil {
		return err
	}

	return nil
}
