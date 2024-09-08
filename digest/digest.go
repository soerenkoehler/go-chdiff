package digest

import (
	"github.com/soerenkoehler/go-chdiff/util"
)

func (digest *Digest) AddFileHash(file, hash string) {
	newHashType := getHashType(hash)
	if newHashType == Unknown {
		util.Fatal("unknown hash type: %v", hash)
	}

	if digest.Algorithm != newHashType {
		if digest.Algorithm != Unknown {
			util.Fatal("hash type mismatch old=%v new=%v", digest.Algorithm, newHashType)
		}
		digest.Algorithm = newHashType
	}

	(*digest.Entries)[file] = hash
}

func getHashType(hash string) HashType {
	switch len(hash) {
	case 128:
		return SHA512
	case 64:
		return SHA256
	default:
		return Unknown
	}
}
