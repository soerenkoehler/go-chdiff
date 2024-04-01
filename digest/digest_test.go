package digest_test

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"github.com/soerenkoehler/go-chdiff/digest"
	"github.com/stretchr/testify/assert"
)

func TestInitialHashUnknown(t *testing.T) {
	d := digest.NewDigest("path", time.Now())
	assert.Equal(t, digest.Unknown, d.Algorithm)
}

func TestHashTypeSetOnFirstCall(t *testing.T) {
	d := digest.NewDigest("path", time.Now())
	d.AddFileHash("file", createRandomHash(32))
	assert.Equal(t, digest.SHA256, d.Algorithm)
}

func TestInvalidHash(t *testing.T) {
	assert.PanicsWithError(t, `/!\ unknown hash type: bad-hash`, func() {
		d := digest.NewDigest("path", time.Now())
		d.AddFileHash("file", "bad-hash")
	})
}

func TestHashTypeMismatch(t *testing.T) {
	assert.PanicsWithError(t, `/!\ hash type mismatch old=1 new=2`, func() {
		d := digest.NewDigest("path", time.Now())
		d.AddFileHash("file", createRandomHash(32))
		d.AddFileHash("file", createRandomHash(64))
	})
}

func createRandomHash(hashsize int) string {
	hashbytes := make([]byte, hashsize)
	rand.Read(hashbytes)
	return hex.EncodeToString(hashbytes)
}
