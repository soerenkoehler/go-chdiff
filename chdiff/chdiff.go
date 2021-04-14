package chdiff

import (
	_ "embed"
	"log"
	"os"
	"path"

	"github.com/alecthomas/kong"
	"github.com/soerenkoehler/chdiff-go/digest"
	"github.com/soerenkoehler/chdiff-go/util"
)

//go:embed description.txt
var _Description string

var cli struct {
	Create cmdDigest `cmd:"" name:"create" aliases:"c" help:"Create digest file for PATH."`
	Verify cmdDigest `cmd:"" name:"verify" aliases:"v" default:"1" help:"Verify digest file for PATH."`
}

type cmdDigest struct {
	Path      string `arg:"" name:"PATH" type:"path" default:"." help:"Path for which to calculate the digest"`
	Algorithm string `name:"alg" help:"The checksum algorithm to use [SHA256,SHA512]." enum:"SHA256,SHA512" default:"SHA256"`
}

func DoMain(
	version string,
	args []string,
	digestService digest.Service,
	stdioService util.StdIOService) {

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
		digestService.Create(
			cli.Create.Path,
			path.Join(cli.Create.Path, "out.txt"),
			cli.Create.Algorithm)

	case "verify", "verify <PATH>":
		digestService.Verify(
			cli.Verify.Path,
			path.Join(cli.Verify.Path, "out.txt"),
			cli.Verify.Algorithm)

	default:
		log.Fatalf("unknown command: %s", ctx.Command())
	}
}
