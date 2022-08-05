package digest

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/soerenkoehler/go-chdiff/util"
)

const SEPARATOR_TEXT = "  "
const SEPARATOR_BINARY = " *"

type Reader func(string, string) (Digest, error)

type Writer func(Digest) error

func Load(path, algorithm string) (Digest, error) {
	if util.Stat(path).IsDir {
		return load(defaultDigestFile(path, algorithm))
	}
	return load(path)
}

func load(digestFile string) (Digest, error) {
	digest := Digest{
		Entries: &FileHashes{}}
	input, err := os.Open(digestFile)
	if err == nil {
		lines := bufio.NewScanner(input)
		for lines.Scan() {
			normalized := strings.Replace(lines.Text(), SEPARATOR_TEXT, SEPARATOR_BINARY, 1)
			tokens := strings.SplitN(normalized, SEPARATOR_BINARY, 2)
			(*digest.Entries)[tokens[1]] = tokens[0]
		}
	}
	return digest, err
}

func Save(digest Digest) error {
	return save(
		defaultDigestFile(digest.Location.Path, digest.Algorithm),
		digest)
}

func save(file string, digest Digest) error {
	// TODO save digest data
	return nil
}

func defaultDigestFile(path, algorithm string) string {
	return fmt.Sprintf("%v/.chdiff.%v.txt", path, algorithm)
}

// TODO
// func (digest Digest) sortedKeys() []string {
// 	keys := make([]string, 0, len(digest))
// 	for key := range digest {
// 		keys = append(keys, key)
// 	}
// 	sort.Strings(keys)
// 	return keys
// }

// TODO
// func (entry DigestEntry) entryToString() string {
// 	return fmt.Sprintf(
// 		"# %d %s %s\n%s *%s\n",
// 		entry.size,
// 		entry.modTime.Local().Format("20060102-150405"),
// 		entry.file,
// 		entry.hash,
// 		entry.file)
// }
