package main

import (
	"flag"
	"fmt"
	"os"
)

const usage1 string = `Usage: %s <COMMAND> [OPTIONS]
Options:
`

const usage2 string = `
Commands:
	keygen genkey                  Generate an Ed25519 publishing keypair.
	keygen dist                    Publish a new release artifact.
	keygen help                    Print help.
	keygen version                 Print keygen version.

Examples:
	keygen genkey
	keygen dist ./dist/app-1-0-0 \
		--account ... \
		--product ... \
		--token ... \
		...

`

type Options struct {
	command  string
	args     []string
	account  string
	product  string
	token    string
	config   string
	logto    string
	loglevel string
}

func ParseArgs() (*Options, error) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage1, os.Args[0])
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, usage2)
	}

	account := flag.String(
		"account",
		"",
		"Your keygen.sh account ID.",
	)

	product := flag.String(
		"product",
		"",
		"Your keygen.sh product ID.",
	)

	token := flag.String(
		"token",
		"",
		"Your keygen.sh product token.",
	)

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
		account:  *account,
		product:  *product,
		token:    *token,
		config:   *config,
		logto:    *logto,
		loglevel: *loglevel,
	}

	switch opts.command {
	case "genkey":
		opts.args = flag.Args()[1:]
	case "dist":
		opts.args = flag.Args()[1:]
	case "version":
		fmt.Println(version)
		os.Exit(0)
	case "help":
		flag.Usage()
		os.Exit(0)
	default:
		flag.Usage()
		os.Exit(1)
	}

	return opts, nil
}
