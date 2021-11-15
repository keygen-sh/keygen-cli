package cmd

import (
	"fmt"

	"github.com/keygen-sh/keygen-cli/internal/spinnerext"
	"github.com/keygen-sh/keygen-go"
	"github.com/mattn/go-tty"
	"github.com/spf13/cobra"
)

type KeyboardKey rune

var (
	KeyboardKeyEnter KeyboardKey = 13
	KeyboardKeyY     KeyboardKey = 121
)

var (
	upgradeCmd = &cobra.Command{
		Use:   "upgrade",
		Short: "check if a CLI upgrade is available",
		Args:  cobra.NoArgs,
		RunE:  upgradeRun,

		SilenceUsage: true,
		Hidden:       true,
	}
)

func init() {
	keygen.UpgradeKey = "5ec69b78d4b5d4b624699cef5faf3347dc4b06bb807ed4a2c6740129f1db7159"
	keygen.PublicKey = "e8601e48b69383ba520245fd07971e983d06d22c4257cfd82304601479cee788"
	keygen.Account = "1fddcec8-8dd3-4d8d-9b16-215cac0f9b52"
	keygen.Product = "2313b7e7-1ea6-4a01-901e-2931de6bb1e2"

	rootCmd.AddCommand(upgradeCmd)
}

func upgradeRun(cmd *cobra.Command, args []string) error {
	tty, err := tty.Open()
	if err != nil {
		// Skip for non-TTYs
		return nil
	}
	defer tty.Close()

	if cmd != nil {
		spinnerext.Start()
	}

	spinnerext.Text("checking for upgrade...")

	release, err := keygen.Upgrade(Version)
	switch {
	case err == keygen.ErrUpgradeNotAvailable:
		if cmd != nil {
			spinnerext.Stop("all up to date!")
		}

		return nil
	case err != nil:
		return err
	}

	spinnerext.Pause()

	fmt.Printf("an upgrade is available! would you like to install v" + release.Version + " now? Y/n ")

	r, err := tty.ReadRune()
	if err != nil {
		return err
	}

	spinnerext.Unpause()

	if k := KeyboardKey(r); k != KeyboardKeyEnter && k != KeyboardKeyY {
		if cmd != nil {
			spinnerext.Stop("upgrade aborted")
		}

		return nil
	}

	spinnerext.Text("installing upgrade...")

	if err := release.Install(); err != nil {
		return err
	}

	if cmd != nil {
		spinnerext.Stop("install complete! now on v" + release.Version)
	}

	return nil
}
