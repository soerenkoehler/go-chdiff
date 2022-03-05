package digest

import (
	"fmt"
)

type Reader func(string) (Digest, error)

type Writer func(Digest) error

func Load(digestFile string) (Digest, error) {
	// TODO
	return Digest{}, nil
}

func Save(digest Digest) error {
	return save(DefaultDigestFile(digest.Location.Path, digest.Algorithm), digest)
}

func DefaultDigestFile(path, algorithm string) string {
	return fmt.Sprintf("%v/.chdiff.%v.txt", path, algorithm)
}

func save(file string, digest Digest) error {
	// TODO
	return nil
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
