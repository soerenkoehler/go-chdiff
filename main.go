package main

import (
	_ "embed"
	"os"

	"github.com/soerenkoehler/go-chdiff/chdiff"
	"github.com/soerenkoehler/go-chdiff/diff"
	"github.com/soerenkoehler/go-chdiff/digest"
)

var _Version = "DEV"

func main() {
	chdiff.Chdiff(
		_Version,
		os.Args,
		chdiff.ChdiffDependencies{
			DigestRead:      digest.Load,
			DigestWrite:     digest.Save,
			DigestCalculate: digest.Calculate,
			DigestCompare:   diff.Compare,
			DiffPrint:       diff.Print,
			Stdout:          os.Stdout,
			Stderr:          os.Stderr,
			KongExit:        os.Exit})
}
