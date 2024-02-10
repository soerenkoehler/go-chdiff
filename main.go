package main

import (
	_ "embed"
	"io"
	"os"

	"github.com/soerenkoehler/go-chdiff/chdiff"
	"github.com/soerenkoehler/go-chdiff/diff"
	"github.com/soerenkoehler/go-chdiff/digest"
)

var _Version = "DEV"

type DefaultDependencies struct{}

func (*DefaultDependencies) DigestRead(dp, df string) (digest.Digest, error) {
	return digest.Load(dp, df)
}

func (*DefaultDependencies) DigestWrite(d digest.Digest, df string) error {
	return digest.Save(d, df)
}

func (*DefaultDependencies) DigestCalculate(rp string, ht digest.HashType) digest.Digest {
	return digest.Calculate(rp, ht)
}

func (*DefaultDependencies) DigestCompare(old, new digest.Digest) diff.Diff {
	return diff.Compare(old, new)
}

func (*DefaultDependencies) DiffPrint(out io.Writer, d diff.Diff) {
	diff.Print(out, d)
}

func (*DefaultDependencies) Stdout() io.Writer {
	return os.Stdout
}

func (*DefaultDependencies) Stderr() io.Writer {
	return os.Stderr
}

func (*DefaultDependencies) KongExit() func(int) {
	return os.Exit
}

func main() {
	chdiff.Chdiff(_Version, os.Args, &DefaultDependencies{})
}
