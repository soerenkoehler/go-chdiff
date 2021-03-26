package main

import (
	_ "embed"
	"os"

	"github.com/soerenkoehler/chdiff-go/chdiff"
	"github.com/soerenkoehler/chdiff-go/digest"
	"github.com/soerenkoehler/chdiff-go/util"
)

var _Version = "DEV"

func main() {
	chdiff.DoMain(
		_Version,
		os.Args,
		digest.DefaultService{},
		util.DefaultStdIOService{})
}
