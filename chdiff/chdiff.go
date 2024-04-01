package chdiff

import (
	_ "embed"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/soerenkoehler/go-chdiff/common"
	"github.com/soerenkoehler/go-chdiff/diff"
	"github.com/soerenkoehler/go-chdiff/digest"
	"github.com/soerenkoehler/go-chdiff/util"
)

const (
	DefaultDigestName  string = ".chdiff.txt"
	UserConfigFileName string = ".chdiff-config.json"
)

var (
	//go:embed description.txt
	_description string

	//go:embed default-config.json
	_defaultConfigJson []byte
)

type cmdDigest struct {
	RootPath   string `arg:"" name:"PATH" type:"path" default:"." help:"Path for which to calculate the digest"`
	DigestFile string `name:"file" short:"f" help:"Optional: Path to different location of the digest file."`
}

type CmdCreate struct {
	cmdDigest
	Algorithm string `name:"algorithm" short:"a" help:"The checksum algorithm to use [SHA256,SHA512]." enum:"SHA256,SHA512" default:"SHA256"`
}

type CmdVerify struct{ cmdDigest }

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
	util.InitLogger(deps.Stderr())

	loadConfig()

	var cli struct {
		Create CmdCreate `cmd:"" name:"create" aliases:"c" help:"Create digest file for PATH."`
		Verify CmdVerify `cmd:"" name:"verify" aliases:"v" help:"Verify digest file for PATH."`
	}

	ctx := kong.Parse(
		&cli,
		kong.Vars{"VERSION": version},
		kong.Description(_description),
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
	var chain util.ChainContext
	var oldDigest digest.Digest

	chain.Chain(func() {
		oldDigest, chain.Err = deps.DigestRead(
			cmd.RootPath,
			defaultDigestFile(cmd.cmdDigest))
	}).Chain(func() {
		deps.DiffPrint(
			deps.Stdout(),
			deps.DigestCompare(
				oldDigest,
				deps.DigestCalculate(
					cmd.RootPath,
					oldDigest.Algorithm)))
	}).ChainError("verify")

	return chain.Err
}

func loadConfig() {
	chain := &util.ChainContext{}
	chain.Chain(func() {
		chain.Err = json.Unmarshal(readConfigFile(), &common.Config)
	}).Chain(func() {
		util.SetLogLevelByName(common.Config.LogLevel)
		util.Debug("%+v", common.Config)
	}).ChainFatal("reading config")
}

func readConfigFile() []byte {
	userhome, err := os.UserHomeDir()
	if err != nil {
		util.Warn("can't determine user home")
		return _defaultConfigJson
	}

	configFile := filepath.Join(userhome, UserConfigFileName)
	data, err := os.ReadFile(configFile)
	if err != nil {
		os.WriteFile(configFile, _defaultConfigJson, 0744)
		return _defaultConfigJson
	}

	return data
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
