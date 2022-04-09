package digest

import (
	"time"

	"github.com/soerenkoehler/go-chdiff/common"
)

type FileHashes map[string]string

type Digest struct {
	Location  common.Location
	Algorithm string
	Entries   *FileHashes
}

func NewDigest(
	path, algorithm string,
	time time.Time) Digest {

	return Digest{
		Location: common.Location{
			Path: path,
			Time: time},
		Algorithm: algorithm,
		Entries:   &FileHashes{}}
}
