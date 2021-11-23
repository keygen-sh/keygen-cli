package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	helpCmd = &cobra.Command{
		Use:   "help [command]",
		Short: "help for a command",
		Run: func(c *cobra.Command, args []string) {
			cmd, _, e := c.Root().Find(args)

			if cmd == nil || e != nil {
				c.Printf("unknown help topic %#q\n", args)

				if err := c.Root().Usage(); err != nil {
					fmt.Fprintln(os.Stderr, "error:", err)
					os.Exit(1)
				}
			} else {
				cmd.InitDefaultHelpFlag()

				if err := cmd.Help(); err != nil {
					fmt.Fprintln(os.Stderr, "error:", err)
					os.Exit(1)
				}
			}
		},
	}
)
