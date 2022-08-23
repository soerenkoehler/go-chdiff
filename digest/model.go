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
	digestPath, algorithm string,
	digestTime time.Time) Digest {

	return Digest{
		Location: common.Location{
			Path: digestPath,
			Time: time.Date(
				digestTime.Local().Year(),
				digestTime.Local().Month(),
				digestTime.Local().Day(),
				digestTime.Local().Hour(),
				digestTime.Local().Minute(),
				digestTime.Local().Second(),
				0,
				time.Local)},
		Algorithm: algorithm,
		Entries:   &FileHashes{}}
}
