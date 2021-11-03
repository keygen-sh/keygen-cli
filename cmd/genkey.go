package cmd

import (
	"fmt"

	"github.com/keygen-sh/keygen-cli/internal/ed25519ph"
	"github.com/spf13/cobra"
)

var (
	genkeyCmd = &cobra.Command{
		Use:   "genkey",
		Short: "generate an Ed25519 keypair for code signing",
		RunE:  genkeyRun,
	}
)

func init() {
	rootCmd.AddCommand(genkeyCmd)
}

func genkeyRun(cmd *cobra.Command, args []string) error {
	pub, priv, err := ed25519ph.GenerateKey()
	if err != nil {
		return err
	}

	// TODO(ezekg) Save to files ./keygen.key and ./keygen.pub
	fmt.Printf("private key:\n  %x\n", priv)
	fmt.Printf("public key:\n  %x\n", pub)

	return nil
}
