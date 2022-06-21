package digest

import (
	"fmt"
	"os"
)

type Reader func(string, string) (Digest, error)

type Writer func(Digest) error

func Load(path, algorithm string) (Digest, error) {
	return load(defaultDigestFile(path, algorithm))
}

func load(digestFile string) (Digest, error) {
	digest := Digest{}
	_, err := os.Open(digestFile)
	// TODO load digest data
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
