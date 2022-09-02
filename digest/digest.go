package digest

import "fmt"

func (digest *Digest) AddFileHash(file, hash string) {
	newHashType := getHashType(hash)
	if digest.Algorithm == Unknown {
		digest.Algorithm = newHashType
	} else if digest.Algorithm != newHashType {
		panic(fmt.Errorf("hash type mismatch old=%v new=%v", newHashType, hash))
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
