package main

import (
	_ "embed"
	"os"

	"github.com/soerenkoehler/chdiff-go/main/chdiff"
	"github.com/soerenkoehler/chdiff-go/main/digest"
	"github.com/soerenkoehler/chdiff-go/main/util"
)

var _Version = "DEV"

func main() {
	chdiff.DoMain(
		_Version,
		os.Args,
		digest.DefaultService{},
		util.DefaultStdIOService{})
}
