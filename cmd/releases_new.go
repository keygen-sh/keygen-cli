package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	releasesNewCmd = &cobra.Command{
		Use:   "new <path>",
		Short: "Publish a new release for a product",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("<path> to file is required")
			}

			path := args[0]
			info, err := os.Stat(path)
			if err != nil {
				reason, ok := err.(*os.PathError)
				if !ok {
					return err
				}

				return fmt.Errorf("<path> is not accessible (%s)", reason.Err)
			}

			if info.IsDir() {
				return errors.New("<path> is a directory (must be a file)")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			ext := filepath.Ext(path)

			fmt.Println("publishing new release...")
			fmt.Printf("  path=%s ext=%s", path, ext)
		},
	}
)

func init() {
	releasesCmd.AddCommand(releasesNewCmd)
}
