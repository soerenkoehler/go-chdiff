package digest

import "fmt"

func (digest *Digest) AddFileHash(file, hash string) {
	newHashType := getHashType(hash)
	if digest.Algorithm != newHashType {
		if digest.Algorithm != Unknown {
			panic(fmt.Errorf("hash type mismatch old=%v new=%v", digest.Algorithm, newHashType))
		}
		digest.Algorithm = newHashType
	}
	(*digest.Entries)[file] = hash
}

func getHashType(hash string) HashType {
	switch len(hash) {
	case 64:
		return SHA256
	case 128:
		return SHA512
	}
	panic(fmt.Errorf("invalid hash %v", hash))
}
