package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	releasesNewCmd = &cobra.Command{
		Use:   "new",
		Short: "Publish a new release for a product",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("publishing new release...")
		},
	}
)

func init() {
	releasesCmd.AddCommand(releasesNewCmd)
}
