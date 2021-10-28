package main

import (
	"flag"
	"fmt"
	"os"
)

const usage1 string = `Usage: %s <command> [subcommand] [options]
Options:
`

const usage2 string = `
Commands:
  keygen help
        Print a useful help message.
  keygen version
        Print keygen version.
  keygen genkey
        Generate an Ed25519 publishing keypair.
  keygen releases <subcommand>
        Manage and publish releases.
        Subcommands:
          new          Publish a new release.
          yank         Yank a release.
          del          Delete a release.
          list         List releases.

Examples:
  keygen genkey
  keygen releases new ./dist/app-1-0-0 \
    --account ... \
    --product ... \
    --token ... \
    ...

`

type Options struct {
	command    string
	subcommand string
	args       []string
	account    string
	product    string
	token      string
	config     string
	logto      string
	loglevel   string
}

func ParseArgs() (*Options, error) {
	shared := flag.NewFlagSet("shared", flag.ExitOnError)
	shared.Usage = func() {
		fmt.Fprintf(os.Stderr, usage1, os.Args[0])
		shared.PrintDefaults()
		fmt.Fprint(os.Stderr, usage2)
	}

	config := shared.String(
		"config",
		"",
		"Path to keygen configuration file. (default: $HOME/.keygen)",
	)

	logto := shared.String(
		"log",
		"none",
		"Write log messages to this file. 'stdout' and 'none' have special meanings.",
	)

	loglevel := shared.String(
		"log-level",
		"DEBUG",
		"The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR.",
	)

	flag.Parse()
	shared.Parse(flag.Args()[:])

	opts := &Options{
		command:  flag.Arg(0),
		config:   *config,
		logto:    *logto,
		loglevel: *loglevel,
	}

	switch opts.command {
	case "genkey":
		opts.args = flag.Args()[1:]
	case "releases":
		cmd := flag.NewFlagSet("releases", flag.ExitOnError)
		cmd.Usage = func() {
			fmt.Fprintf(os.Stderr, usage1, os.Args[0])
			cmd.PrintDefaults()
			fmt.Fprint(os.Stderr, usage2)
		}

		opts.subcommand = flag.Arg(1)

		account := cmd.String(
			"account",
			"",
			"Your keygen.sh account ID.",
		)

		product := cmd.String(
			"product",
			"",
			"Your keygen.sh product ID.",
		)

		token := cmd.String(
			"token",
			"",
			"Your keygen.sh product token.",
		)

		cmd.Parse(flag.Args()[1:])

		switch opts.subcommand {
		case "":
			fmt.Println("Error: you must provide a subcommand")
			cmd.Usage()
			os.Exit(1)
		case "help":
			cmd.Usage()
			os.Exit(0)
		}

		opts.args = flag.Args()[2:]
		opts.account = *account
		opts.product = *product
		opts.token = *token
	case "version":
		fmt.Println(version)
		os.Exit(0)
	case "help":
		shared.Usage()
		os.Exit(0)
	default:
		fmt.Println("Error: you must provide a valid command")
		shared.Usage()
		os.Exit(1)
	}

	return opts, nil
}
