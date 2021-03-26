package chdiff

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/soerenkoehler/chdiff-go/main/digest"
	"github.com/soerenkoehler/chdiff-go/main/util"
)

//go:embed description.txt
var _Description string

var cli struct {
	Create cmdDigest `cmd:"" name:"create" aliases:"c" help:"Create digest file for PATH."`
	Verify cmdDigest `cmd:"" name:"verify" aliases:"v" default:"1" help:"Verify digest file for PATH."`
}

type cmdDigest struct {
	Path string `arg:"" name:"PATH" type:"path" default:"." help:"Path for which to calculate the digest"`
	Mode string `name:"mode" short:"m" help:"The checksum algorithm to use [SHA256,SHA512]." enum:"SHA256,SHA512" default:"SHA256"`
}

func DoMain(
	version string,
	args []string,
	digestService digest.Service,
	stdioService util.StdIOService) {

	var err error

	os.Args = args
	log.SetOutput(stdioService.Stdout())

	ctx := kong.Parse(
		&cli,
		kong.Vars{"VERSION": version},
		kong.Description(_Description),
		kong.UsageOnError(),
		kong.Writers(
			stdioService.Stdout(),
			stdioService.Stderr()))

	switch ctx.Command() {

	case "create", "create <PATH>":
		err = digestService.Create(cli.Create.Path, "out.txt", cli.Create.Mode)

	case "verify", "verify <PATH>":
		err = digestService.Verify(cli.Verify.Path, "out.txt", cli.Verify.Mode)

	default:
		err = fmt.Errorf("unknown command: %s", ctx.Command())
	}

	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
