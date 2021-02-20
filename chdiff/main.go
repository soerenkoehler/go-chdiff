package main

import (
	"fmt"
	"log"

	"github.com/alecthomas/kong"
	"github.com/soerenkoehler/chdiff-go/digest"
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
	var err error

	ctx := kong.Parse(&cli, kong.UsageOnError(), kong.Description("TODO: Description"))

	switch ctx.Command() {

	case "create", "create <PATH>":
		err = digest.Create(cli.Create.Path, "out.txt", cli.Create.Mode)

	case "verify", "verify <PATH>":
		err = digest.Verify(cli.Verify.Path, "out.txt", cli.Verify.Mode)

	case "version":
		log.Printf("Version: %s", _Version)

	default:
		err = fmt.Errorf("unknown command: %s", ctx.Command())
	}

	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
