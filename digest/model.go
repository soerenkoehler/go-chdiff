package digest

import (
	"time"

	"github.com/soerenkoehler/go-chdiff/common"
)

type digestEntry struct {
	file    string
	hash    string
	size    int64
	modTime time.Time
}

type digestEntries map[string]digestEntry

type Digest struct {
	Location common.Location
	Entries  *digestEntries
}

func NewDigest(
	path string,
	time time.Time) Digest {

	return Digest{
		Location: common.Location{
			Path: path,
			Time: time},
		Entries: &digestEntries{}}
}

func (digest Digest) AddNewEntry(
	file string,
	hash string,
	size int64,
	time time.Time) Digest {

	return digest.AddEntry(digestEntry{
		file:    file,
		hash:    hash,
		size:    size,
		modTime: time})
}

func (digest Digest) AddEntry(entry digestEntry) Digest {

	(*digest.Entries)[entry.file] = entry

	return digest
}
