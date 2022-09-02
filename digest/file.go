package digest

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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

	digest := NewDigest(digestPath, "todo", digestFileInfo.ModTime().Local())

	input, err := os.Open(digestFile)
	if err == nil {
		lines := bufio.NewScanner(input)
		for lines.Scan() {
			normalized := strings.Replace(lines.Text(), SEPARATOR_TEXT, SEPARATOR_BINARY, 1)
			tokens := strings.SplitN(normalized, SEPARATOR_BINARY, 2)
			if len(tokens) != 2 {
				return Digest{}, fmt.Errorf("invalid digest file")
			}
			(*digest.Entries)[tokens[1]] = tokens[0]
		}
	}
	return digest, err
}

// func Save(digest Digest) error {
// 	return save(defaultDigestFile(digest.Location.Path), digest)
// }

func Save(digest Digest, digestFile string) error {
	// TODO save digest data
	output, err := os.Create(digestFile)
	if err == nil {
		for k, v := range *digest.Entries {
			fmt.Fprintf(output, "%v *%v", k, v)
		}
		os.Chtimes(digestFile, digest.Location.Time, digest.Location.Time)
	}
	return nil
}

func defaultDigestFile(path string) string {
	return filepath.Join(path, ".chdiff.txt")
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
