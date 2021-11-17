package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/keygen-sh/keygen-go"
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
	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return nil
	}

	release, err := keygen.Upgrade(Version)
	switch {
	case err == keygen.ErrUpgradeNotAvailable:
		if cmd != nil {
			fmt.Println("all up to date!")
		}

		return nil
	case err != nil:
		return err
	}

	fmt.Printf("an upgrade is available! would you like to install v" + release.Version + " now? Y/n ")

	key, _, err := keyboard.GetSingleKey()
	if err != nil {
		return err
	}

	fmt.Println()

	if k := KeyCode(key); k != KeyCodeEnter && k != KeyCodeY {
		if cmd != nil {
			fmt.Println("upgrade aborted")
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
		fmt.Println("install complete! now on v" + release.Version)
	}

	return nil
}
