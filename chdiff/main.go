package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/soerenkoehler/chdiff-go/digest"
	"github.com/soerenkoehler/chdiff-go/util"
)

var _Version = "DEV"

var cli struct {
	Version cmdVersion `cmd:"" name:"version" help:"Show detailed version info."`
	Create  cmdDigest  `cmd:"" name:"create" aliases:"c" help:"Create digest file for PATH."`
	Verify  cmdDigest  `cmd:"" name:"verify" aliases:"v" default:"1" help:"Verify digest file for PATH."`
}

type cmdVersion struct {
}

type cmdDigest struct {
	Path string `arg:"" name:"PATH" type:"path" default:"." help:"Path for which to calculate the digest"`
	Mode string `name:"mode" short:"m" help:"The checksum algorithm to use [SHA256,SHA512]." enum:"SHA256,SHA512" default:"SHA256"`
}

func main() {
	if errors := doMain(); len(errors) > 0 {
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "Error: %v\n", e)
		}
		os.Exit(1)
	}
}

func doMain() util.ErrorList {
	ctx := kong.Parse(&cli, kong.UsageOnError(), kong.Description("TODO: Description"))

	switch ctx.Command() {

	case "create", "create <PATH>":
		return digest.Create(cli.Create.Path, "out.txt", cli.Create.Mode)

	case "verify", "verify <PATH>":
		return digest.Verify(cli.Verify.Path, "out.txt", cli.Verify.Mode)

	case "version":
		fmt.Fprintln(os.Stdout, _Version)
		return util.ErrorList{}

	default:
		panic(ctx.Command())
	}
}
