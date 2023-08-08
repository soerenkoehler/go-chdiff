package digest

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const SEPARATOR_TEXT = "  "
const SEPARATOR_BINARY = " *"

type Reader func(digestRootPath, digestFile string) (Digest, error)

type Writer func(digest Digest, digestFile string) error

func Load(digestPath, digestFile string) (Digest, error) {
	digestFileInfo, err := os.Lstat(digestFile)
	if err != nil {
		return Digest{}, err
	}

	digest := NewDigest(digestPath, digestFileInfo.ModTime().Local())

	input, err := os.Open(digestFile)
	if err == nil {
		lines := bufio.NewScanner(input)
		for lines.Scan() {
			normalized := strings.Replace(lines.Text(), SEPARATOR_TEXT, SEPARATOR_BINARY, 1)
			println(normalized)
			tokens := strings.SplitN(normalized, SEPARATOR_BINARY, 2)
			if len(tokens) != 2 {
				return Digest{}, fmt.Errorf("invalid digest file")
			}
			digest.AddFileHash(tokens[1], tokens[0])
		}
	}
	return digest, err
}

func Save(digest Digest, digestFile string) error {
	// TODO save digest data
	output, err := os.Create(digestFile)
	if err == nil {
		for k, v := range *digest.Entries {
			fmt.Fprintf(output, "%v%v%v\n", v, SEPARATOR_BINARY, k)
		}
		os.Chtimes(digestFile, digest.Location.Time, digest.Location.Time)
	}
	return nil
}
