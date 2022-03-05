package digest

import (
	"time"

	"github.com/soerenkoehler/go-chdiff/common"
)

type digestEntry struct {
	file string
	Hash string
}

type digestEntries map[string]digestEntry

type Digest struct {
	Location  common.Location
	Algorithm string
	Entries   *digestEntries
}

func newDigest(
	path, algorithm string,
	time time.Time) Digest {

	return Digest{
		Location: common.Location{
			Path: path,
			Time: time},
		Algorithm: algorithm,
		Entries:   &digestEntries{}}
}

func (digest Digest) addEntry(entry digestEntry) Digest {

	(*digest.Entries)[entry.file] = entry

	return digest
}
