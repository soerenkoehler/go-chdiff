package chdiff

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/soerenkoehler/go-chdiff/diff"
	"github.com/soerenkoehler/go-chdiff/digest"
)

const DefaultDigestName string = ".chdiff.txt"

//go:embed description.txt
var _Description string

var cli struct {
	Create CmdCreate `cmd:"" name:"create" aliases:"c" help:"Create digest file for PATH."`
	Verify CmdVerify `cmd:"" name:"verify" aliases:"v" help:"Verify digest file for PATH."`
}

type CmdCreate struct {
	cmdDigest
	Algorithm string `name:"algorithm" short:"a" help:"The checksum algorithm to use [SHA256,SHA512]." enum:"SHA256,SHA512" default:"SHA256"`
}

type CmdVerify struct{ cmdDigest }

type cmdDigest struct {
	RootPath   string `arg:"" name:"PATH" type:"path" default:"." help:"Path for which to calculate the digest"`
	DigestFile string `name:"file" short:"f" help:"Optional: Path to different location of the digest file."`
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
			cmd.RootPath,
			hashTypeFromAlgorithm(cmd.Algorithm)),
		defaultDigestFile(cmd.cmdDigest))
}

func (cmd *CmdVerify) Run(deps ChdiffDependencies) error {
	oldDigest, err := deps.DigestRead(
		cmd.RootPath,
		defaultDigestFile(cmd.cmdDigest))
	if err == nil {
		deps.DiffPrint(
			deps.Stdout,
			deps.DigestCompare(
				oldDigest,
				deps.DigestCalculate(
					cmd.RootPath,
					oldDigest.Algorithm)))
	}
	return err
}

func hashTypeFromAlgorithm(algorithm string) digest.HashType {
	switch algorithm {
	case "SHA256":
		return digest.SHA256
	case "SHA512":
		return digest.SHA512
	}
	panic(fmt.Errorf("invalid algorithm %v", algorithm))
}

func defaultDigestFile(cmd cmdDigest) string {
	if len(cmd.DigestFile) > 0 {
		return cmd.DigestFile
	}
	return filepath.Join(cmd.RootPath, DefaultDigestName)
}
