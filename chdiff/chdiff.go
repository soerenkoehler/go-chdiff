package chdiff

import (
	_ "embed"
	"io"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/soerenkoehler/go-chdiff/diff"
	"github.com/soerenkoehler/go-chdiff/digest"
)

//go:embed description.txt
var _Description string

var cli struct {
	Create CmdCreate `cmd:"" name:"create" aliases:"c" help:"Create digest file for PATH."`
	Verify CmdVerify `cmd:"" name:"verify" aliases:"v" help:"Verify digest file for PATH."`
}

type CmdCreate struct{ cmdDigest }

type CmdVerify struct{ cmdDigest }

type cmdDigest struct {
	rootPath   string `arg:"" name:"PATH" type:"path" default:"." help:"Path for which to calculate the digest"`
	digestFile string `name:"file" aliases:"f" help:"Optional: Path to different location of the digest file."`
	Algorithm  string `name:"algorithm" aliases:"a,algo" help:"The checksum algorithm to use [SHA256,SHA512]." enum:"SHA256,SHA512" default:"SHA256"`
}

type ChdiffDependencies struct {
	DigestRead      digest.Reader
	DigestWrite     digest.Writer
	DigestCalculate digest.Calculator
	DigestCompare   diff.Comparator
	DiffPrint       diff.DiffPrinter
	Stdout          io.Writer
	Stderr          io.Writer
	KongExit        func(int)
}

func Chdiff(
	version string,
	args []string,
	deps ChdiffDependencies) {

	os.Args = args
	log.SetOutput(deps.Stderr)

	ctx := kong.Parse(
		&cli,
		kong.Vars{"VERSION": version},
		kong.Description(_Description),
		kong.Exit(deps.KongExit),
		kong.UsageOnError(),
		kong.Writers(deps.Stdout, deps.Stderr))

	if ctx != nil {
		ctx.FatalIfErrorf(ctx.Run(deps))
	}
}

func (cmd *CmdCreate) Run(deps ChdiffDependencies) error {
	return deps.DigestWrite(
		deps.DigestCalculate(
			cmd.rootPath,
			cmd.Algorithm),
		cmd.digestFile)
}

func (cmd *CmdVerify) Run(deps ChdiffDependencies) error {
	oldDigest, err := deps.DigestRead(
		cmd.rootPath,
		cmd.digestFile)
	if err == nil {
		deps.DiffPrint(
			deps.Stdout,
			deps.DigestCompare(
				oldDigest,
				deps.DigestCalculate(
					cmd.rootPath,
					cmd.Algorithm)))
	}
	return err
}
