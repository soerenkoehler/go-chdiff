package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/soerenkoehler/chdiff/other"
)

const _Version = "ChecksumDiff 0.9"
const _Usage = _Version + `

Usage:
    chdiff (c | create) PATH
    chdiff (v | verify) PATH
    chdiff (-h | --help | --version)

Commands:
    c create  Create checksum file in directory PATH.
    v verify  Verify checksum file in directory PATH.

Options:
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
