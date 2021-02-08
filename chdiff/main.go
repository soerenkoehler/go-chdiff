package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/soerenkoehler/chdiff-go/other"
	"github.com/soerenkoehler/chdiff-go/resource"
	"github.com/soerenkoehler/chdiff-go/util"
)

var _Version = "DEV"

func main() {
	if errors := doMain(); len(errors) > 0 {
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "Error: %v\n", e)
		}
		os.Exit(1)
	}
}

func doMain() []error {
	opts, err := docopt.ParseArgs(
		util.ReplaceVariable(
			resource.Usage,
			"VERSION",
			_Version),
		nil,
		_Version)
	if err == nil {
		other.Func(normalizeOpts(opts))
	}
	return []error{err}
}

// TODO
func normalizeOpts(opts docopt.Opts) docopt.Opts {
	return opts
}
