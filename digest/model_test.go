package digest_test

import (
	"testing"
	"time"

	"github.com/soerenkoehler/go-chdiff/digest"
	"github.com/stretchr/testify/assert"
)

func TestNewDigestTruncatesTime(t *testing.T) {
	digest := digest.NewDigest("",
		time.Date(1999, 12, 31, 23, 59, 58, 999999, time.Local))
	assert.Equal(t, 0, digest.Location.Time.Nanosecond())
}
