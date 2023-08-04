package digest_test

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"github.com/soerenkoehler/go-chdiff/digest"
	"github.com/stretchr/testify/require"
)

func TestInvalidHash(t *testing.T) {
	require.PanicsWithError(t, "invalid hash bad-hash", func() {
		digest := digest.NewDigest("path", time.Now())
		digest.AddFileHash("file", "bad-hash")
	})
}

func TestHashTypeMismatch(t *testing.T) {
	require.PanicsWithError(t, "hash type mismatch old=1 new=2", func() {
		digest := digest.NewDigest("path", time.Now())
		digest.AddFileHash("file", createRandomHash(32))
		digest.AddFileHash("file", createRandomHash(64))
	})
}

func createRandomHash(hashsize int) string {
	hashbytes := make([]byte, hashsize)
	rand.Read(hashbytes)
	return hex.EncodeToString(hashbytes)
}
