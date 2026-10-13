package main

import (
	"flag"
	"os"
)

// flagRefreshData reports whether -refresh-data was passed on the command
// line. Parsed before any other mode flag.
func flagRefreshData() bool {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	refresh := fs.Bool("refresh-data", false, "download the latest API catalog and rewrite gendoc/data/*.json")
	_ = fs.Parse(os.Args[1:])
	return *refresh
}

// flagProbeResource returns the value of -probe=<resource-short-name>.
// Empty when the flag is absent.
func flagProbeResource() string {
	// parse only our flag without disturbing the standard flag set used by
	// potential future CLI options.
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	res := fs.String("probe", "", "resource short name (e.g. vpc) to probe API references")
	_ = fs.Parse(os.Args[1:])
	return *res
}
