package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	helpCmd = &cobra.Command{
		Use:          "help [command]",
		Short:        "help for a command",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			cmd, _, e := c.Root().Find(args)

			if cmd == nil || e != nil {
				fmt.Printf("warning: unknown help topic %#q\n", args)

				return e
			}

			if err := cmd.Help(); err != nil {
				return err
			}

			return nil
		},
	}
)
