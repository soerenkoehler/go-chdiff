package digest_test

import (
	"testing"
	"time"

	"github.com/soerenkoehler/go-chdiff/digest"
)

func TestNewDigestTruncatesTime(t *testing.T) {
	digest := digest.NewDigest("",
		time.Date(1999, 12, 31, 23, 59, 58, 999999, time.Local))
	if digest.Location.Time.Nanosecond() != 0 {
		t.Fatal("NewDigest() should truncate nanoseconds")
	}
}
