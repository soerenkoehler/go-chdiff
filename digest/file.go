package digest

type Reader func(string) (Digest, error)

type Writer func(string, Digest) error

func Load(digestFile string) (Digest, error) {
	// TODO
	return Digest{}, nil
}

func Save(digestFile string, digest Digest) error {
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
