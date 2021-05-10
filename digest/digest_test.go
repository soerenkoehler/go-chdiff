package digest

import (
	"path"
	"testing"

	"github.com/soerenkoehler/go-testutils/datautil"
	// "github.com/google/go-cmp/cmp"
)

func TestDummy(t *testing.T) {
	const fileName = "data"
	const expectedHash = "a6452fbd8c12f8df622c1ca4c567f966801fb56442aca03b4e1303e7a412a9d5"

	tmpdir := t.TempDir()
	tmpfile := path.Join(tmpdir, fileName)

	datautil.CreateRandomFile(tmpfile, 256, 1)

	digest := calculateDigest(tmpdir, "SHA256")

	if digest[fileName].file != fileName {
		t.Error("Digest map key and DigestEntry.file must match")
	}

	if digest[fileName].hash != expectedHash {
		t.Errorf("expected hash: %s actual hash: %s", expectedHash, digest[fileName].hash)
	}
}
