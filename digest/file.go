package digest

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/soerenkoehler/go-chdiff/util"
)

const SEPARATOR_TEXT = "  "
const SEPARATOR_BINARY = " *"

type Reader func(digestRootPath, digestFile string) (Digest, error)

type Writer func(digest Digest, digestFile string) error

func Load(digestPath, digestFile string) (Digest, error) {
	var fileInfo fs.FileInfo
	var input *os.File
	var digest Digest

	chain := &util.ChainContext{}

	chain.Chain(func() {
		fileInfo, chain.Err = os.Lstat(digestFile)
		hackLstatErrorForWindow(&chain.Err)
	}).Chain(func() {
		input, chain.Err = os.Open(digestFile)
	}).Chain(func() {
		defer input.Close()

		digest = NewDigest(digestPath, fileInfo.ModTime().Local())

		lines := bufio.NewScanner(input)
		for lines.Scan() {
			normalized := strings.Replace(lines.Text(), SEPARATOR_TEXT, SEPARATOR_BINARY, 1)
			tokens := strings.SplitN(normalized, SEPARATOR_BINARY, 2)
			if len(tokens) != 2 {
				chain.Err = fmt.Errorf("invalid digest file")
				return
			}
			digest.AddFileHash(tokens[1], tokens[0])
		}
	})

	return digest, chain.Err
}

// hack for better error message under Windows
func hackLstatErrorForWindow(err *error) {
	if *err != nil {
		(*err).(*os.PathError).Op = "lstat"
	}
}

func Save(digest Digest, digestFile string) error {
	var output *os.File

	chain := &util.ChainContext{}

	chain.Chain(func() {
		output, chain.Err = os.Create(digestFile)
	}).Chain(func() {
		defer output.Close()

		for k, v := range *digest.Entries {
			fmt.Fprintf(output, "%v%v%v\n", v, SEPARATOR_BINARY, k)
		}
		os.Chtimes(digestFile, digest.Location.Time, digest.Location.Time)
	})

	return chain.Err
}
