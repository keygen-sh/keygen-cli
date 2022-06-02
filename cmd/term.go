package cmd

import "github.com/fatih/color"

var (
	italic = color.New(color.Italic).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
)
