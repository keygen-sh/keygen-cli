package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/keygen-sh/keygen-go/v2"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

type KeyCode rune

var (
	KeyCodeEnter KeyCode = 13
	KeyCodeY     KeyCode = 121
)

var (
	upgradeOpts = &UpgradeCommandOptions{}
	upgradeCmd  = &cobra.Command{
		Use:   "upgrade",
		Short: "check if a CLI upgrade is available",
		Args:  cobra.NoArgs,
		RunE:  upgradeRun,

		SilenceUsage: true,
		Hidden:       true,
	}
)

type UpgradeCommandOptions struct {
	Channel string
}

func init() {
	upgradeCmd.Flags().StringVar(&upgradeOpts.Channel, "channel", "<auto>", "channel for upgrade check, one of: stable, rc, beta, alpha, dev")

	keygen.PublicKey = "b8f3eb4cd260135f67a5096e8dc1c9b9dcb81ee9fe50d12cdcd941f6607a9031"
	keygen.Account = "5cc3b5a2-0d08-4291-940b-41c21f0ba6ab"
	keygen.Product = "0d5f0b57-3102-4ddf-beb9-f652cf8e24b7"

	rootCmd.AddCommand(upgradeCmd)
}

func upgradeRun(cmd *cobra.Command, args []string) error {
	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return nil
	}

	// When the upgrade command is not called directly, we only want to
	// check periodically. To do so, we'll try to use a /tmp lockfile.
	if cmd == nil {
		path := filepath.Join(os.TempDir(), "keygen-auto-upgrade.lock")
		info, err := os.Stat(path)
		if err != nil {
			f, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
			defer f.Close()
		} else {
			// Check for upgrades at most once per day
			if time.Since(info.ModTime()) < time.Duration(24)*time.Hour {
				return nil
			}

			// Touch the lockfile
			os.Chtimes(path, time.Now(), time.Now())
		}
	}

	// Adjust upgrade channel based on current version
	var channel string
	if upgradeOpts.Channel == "<auto>" {
		switch {
		case strings.Contains(Version, "-rc"):
			channel = "rc"
		case strings.Contains(Version, "-beta"):
			channel = "beta"
		case strings.Contains(Version, "-alpha"):
			channel = "alpha"
		case strings.Contains(Version, "-dev"):
			channel = "dev"
		default:
			channel = "stable"
		}
	} else {
		channel = upgradeOpts.Channel
	}

	release, err := keygen.Upgrade(
		keygen.UpgradeOptions{CurrentVersion: Version, Channel: channel, PublicKey: "5ec69b78d4b5d4b624699cef5faf3347dc4b06bb807ed4a2c6740129f1db7159"},
	)
	switch {
	case err == keygen.ErrUpgradeNotAvailable:
		if cmd != nil {
			fmt.Println("all up to date!")
		}

		return nil
	case err != nil:
		if cmd != nil {
			return err
		}

		return nil
	}

	fmt.Printf("an upgrade is available! would you like to install " + italic("v"+release.Version) + " now? Y/n ")

	key, _, err := keyboard.GetSingleKey()
	if err != nil {
		return err
	}

	fmt.Println()

	if k := KeyCode(key); k != KeyCodeEnter && k != KeyCodeY {
		if cmd != nil {
			fmt.Println(yellow("warning:") + " upgrade aborted")
		}

		return nil
	}

	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	style := mpb.SpinnerStyle(frames...)
	style.PositionLeft()

	progress := mpb.New(mpb.WithWidth(1), mpb.WithRefreshRate(180*time.Millisecond))
	spinner := progress.Add(1,
		mpb.NewBarFiller(style),
		mpb.BarRemoveOnComplete(),
		mpb.AppendDecorators(
			decor.Name("installing..."),
		),
	)

	if err := release.Install(); err != nil {
		return err
	}

	spinner.Increment()
	progress.Wait()

	if cmd != nil {
		fmt.Println(green("upgraded:") + " now on " + italic("v"+release.Version))
	}

	return nil
}
