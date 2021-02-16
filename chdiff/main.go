package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/soerenkoehler/chdiff-go/digest"
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
		return processOpts(opts)
	}
	return []error{err}
}

func processOpts(opts docopt.Opts) []error {
	fmt.Println(opts)
	mode, err := opts.String("MODE")

	if err != nil {
		return []error{err}
	}

	switch {

	case opts["c"]:
		return digest.Create(getPath(opts), "out.txt", mode)

	case opts["v"]:
		return digest.Verify(getPath(opts), "out.txt", mode)

	default:
		fmt.Println("Using default: chdiff v .")
		return digest.Verify(getPath(opts), "out.txt", mode)

	}
}

// func hasOption(options docopt.Opts, name string) bool {
// 	result, err := options.Bool(name)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Checking Option: %v\nError: %v\n", name, err)
// 	}
// 	return result
// }

func getPath(options docopt.Opts) string {
	result, err := options.String("PATH")
	if err == nil {
		return result
	}
	return "."
}
