package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/soerenkoehler/chdiff-go/other"
)

const _Version = "ChecksumDiff 0.9"
const _Usage = _Version + `

Usage:
    chdiff c PATH [-f FILE]
    chdiff v PATH [-f FILE]
    chdiff (-h | --help | --version)

Commands:
    c  Create checksum file in directory PATH.
    v  Verify checksum file in directory PATH.

Options:
	-f FILE    Use the given snapshot file.
    -h --help  Show help.
    --version  Show version.`

func main() {
	opts, err := docopt.ParseArgs(_Usage, nil, _Version)
	if err == nil {
		other.Func(normalizeOpts(opts))
	} else {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(127)
	}
}

// TODO
func normalizeOpts(opts docopt.Opts) docopt.Opts {
	return opts
}
