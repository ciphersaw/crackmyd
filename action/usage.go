package action

import (
	"flag"
	"fmt"
	"os"
)

// usageMap records the usage of each argument.
var usageMap = map[string]string{
	"version":  "Print the version of crackmyd.",
	"password": "Assign the user-defined dictionary of passwords for cracking.",
	"suffix":   "Assign the user-defined dictionary of suffixes for cracking.",
}

// usage customizes the usage information for crackmyd.
func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <file>\n", os.Args[0])
	flag.PrintDefaults()
}
