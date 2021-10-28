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
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage1, os.Args[0])
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, usage2)
	}

	config := flag.String(
		"config",
		"",
		"Path to keygen configuration file. (default: $HOME/.keygen)",
	)

	logto := flag.String(
		"log",
		"none",
		"Write log messages to this file. 'stdout' and 'none' have special meanings.",
	)

	loglevel := flag.String(
		"log-level",
		"DEBUG",
		"The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR.",
	)

	flag.Parse()

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
		if flag.NArg() == 1 {
			fmt.Println("Error: you must provide a subcommand")
			flag.Usage()
			os.Exit(1)
		}

		set := flag.NewFlagSet("releases", flag.ExitOnError)
		set.Usage = func() {
			fmt.Fprintf(os.Stderr, usage1, os.Args[0])
			set.PrintDefaults()
			fmt.Fprint(os.Stderr, usage2)
		}

		subcommand := flag.Arg(1)
		args := flag.Args()[2:]

		account := set.String(
			"account",
			"",
			"Your keygen.sh account ID.",
		)

		product := set.String(
			"product",
			"",
			"Your keygen.sh product ID.",
		)

		token := set.String(
			"token",
			"",
			"Your keygen.sh product token.",
		)

		set.Parse(args)

		if subcommand == "help" {
			set.Usage()
			os.Exit(0)
		}

		opts.subcommand = subcommand
		opts.args = args
		opts.account = *account
		opts.product = *product
		opts.token = *token
	case "version":
		fmt.Println(version)
		os.Exit(0)
	case "help":
		flag.Usage()
		os.Exit(0)
	default:
		fmt.Println("Error: you must provide a command")
		flag.Usage()
		os.Exit(1)
	}

	return opts, nil
}
