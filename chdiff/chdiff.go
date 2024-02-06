package chdiff

import (
	_ "embed"
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

type ChdiffDependencies interface {
	DigestRead(string, string) (digest.Digest, error)
	DigestWrite(digest.Digest, string) error
	DigestCalculate(string, digest.HashType) digest.Digest
	DigestCompare(digest.Digest, digest.Digest) diff.Diff
	DiffPrint(io.Writer, diff.Diff)
	Stdout() io.Writer
	Stderr() io.Writer
	KongExit() func(int)
}

func Chdiff(
	version string,
	args []string,
	deps ChdiffDependencies) {

	os.Args = args
	log.SetOutput(deps.Stderr())

	ctx := kong.Parse(
		&cli,
		kong.Vars{"VERSION": version},
		kong.Description(_Description),
		kong.UsageOnError(),
		kong.Writers(deps.Stdout(), deps.Stderr()),
		kong.Exit(deps.KongExit()),
		kong.BindTo(deps, (*ChdiffDependencies)(nil)))

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
			deps.Stdout(),
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
	case "SHA512":
		return digest.SHA512
	case "SHA256":
		fallthrough
	default:
		return digest.SHA256
	}
}

func defaultDigestFile(cmd cmdDigest) string {
	digestFile := cmd.DigestFile
	if len(cmd.DigestFile) == 0 {
		digestFile = filepath.Join(cmd.RootPath, DefaultDigestName)
	}
	absPath, _ := filepath.Abs(digestFile)
	return absPath
}
