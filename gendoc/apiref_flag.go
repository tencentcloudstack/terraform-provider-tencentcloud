package main

import (
	"flag"
	"os"
)

// modeFlags holds the parsed command-line mode switches.
type modeFlags struct {
	refresh bool
	probe   string
}

// parseModeFlags parses all gendoc mode switches in a single flag set so
// that unknown-flag errors are reported consistently.
func parseModeFlags(args []string) modeFlags {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var m modeFlags
	fs.BoolVar(&m.refresh, "refresh", false,
		"refresh all bundled snapshots (catalog + English descriptions) and exit")
	fs.StringVar(&m.probe, "probe", "",
		"resource short name (e.g. vpc) to probe API references; use 'all' for coverage stats")
	_ = fs.Parse(args)
	return m
}
