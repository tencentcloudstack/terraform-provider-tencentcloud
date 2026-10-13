package main

import (
	"flag"
	"os"
)

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
