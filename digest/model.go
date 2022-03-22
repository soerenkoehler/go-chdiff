package digest

import (
	"time"

	"github.com/soerenkoehler/go-chdiff/common"
)

type DigestEntry struct {
	File string
	Hash string
}

type digestEntries struct {
	common.Set[DigestEntry]
}

type Digest struct {
	Location  common.Location
	Algorithm string
	Entries   *digestEntries
}

func NewDigest(
	path, algorithm string,
	time time.Time) Digest {

	return Digest{
		Location: common.Location{
			Path: path,
			Time: time},
		Algorithm: algorithm,
		Entries:   &digestEntries{}}
}

func (digest Digest) addEntry(entry DigestEntry) Digest {

	digest.Entries.Put(entry)
	(*digest.Entries)[entry.File] = entry

	return digest
}
