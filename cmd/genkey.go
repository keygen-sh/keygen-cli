package cmd

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/keygen-sh/keygen-cli/internal/ed25519ph"
	"github.com/spf13/cobra"
)

var (
	genkeyCmd = &cobra.Command{
		Use:   "genkey",
		Short: "generate an Ed25519 keypair for code signing",
		Args:  genkeyArgs,
		RunE:  genkeyRun,

		// Encountering an error should not display usage
		SilenceUsage: true,
	}
)

func init() {
	genkeyCmd.Flags().StringVar(&flags.privateKey, "private-key", "keygen.key", "write private key to this file")
	genkeyCmd.Flags().StringVar(&flags.publicKey, "public-key", "keygen.pub", "write public key to this file")

	rootCmd.AddCommand(genkeyCmd)
}

func genkeyArgs(cmd *cobra.Command, args []string) error {
	if _, err := os.Stat(flags.privateKey); err == nil {
		return fmt.Errorf(`private key file "%s" already exists`, flags.privateKey)
	}

	if _, err := os.Stat(flags.publicKey); err == nil {
		return fmt.Errorf(`public key file "%s" already exists`, flags.publicKey)
	}

	return nil
}

func genkeyRun(cmd *cobra.Command, args []string) error {
	pub, priv, err := ed25519ph.GenerateKey()
	if err != nil {
		return err
	}

	if err := writePrivateKeyFile(priv); err != nil {
		return err
	}

	if err := writePublicKeyFile(pub); err != nil {
		return err
	}

	return nil
}

func writePrivateKeyFile(privateKey ed25519ph.PrivateKey) error {
	file, err := os.Create("keygen.key")
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

func writePublicKeyFile(publicKey ed25519ph.PublicKey) error {
	file, err := os.Create("keygen.pub")
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
