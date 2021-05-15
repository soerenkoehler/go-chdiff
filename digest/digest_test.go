package digest

import (
	"path"
	"testing"

	"github.com/soerenkoehler/go-testutils/datautil"
	// "github.com/google/go-cmp/cmp"
)

type testCase struct {
	path string
	size int64
	seed uint64
	hash string
}

func runDigestCalculationTest(t *testing.T, data []testCase) {
	root := t.TempDir()

	for _, dataPoint := range data {
		file := path.Join(root, dataPoint.path)
		datautil.CreateRandomFile(file, dataPoint.size, dataPoint.seed)
	}

	digest := calculateDigest(root, "SHA256")

	if len(digest) != len(data) {
		t.Fatal("Digest size must match number of input data points")
	}

	for _, dataPoint := range data {
		expectedPath := dataPoint.path
		actualPath := digest[expectedPath].file
		if actualPath != expectedPath {
			t.Errorf("DigestEntry.file (%v) must match Digest map key (%v)",
				actualPath,
				expectedPath)
		}

		expectedHash := dataPoint.hash
		actualHash := digest[expectedPath].hash
		if actualHash != expectedHash {
			t.Errorf("actual hash (%v) must match expected hash (%s)",
				actualHash,
				expectedHash)
		}
	}
}

func TestDigest(t *testing.T) {
	runDigestCalculationTest(t, []testCase{{
		path: "zero",
		size: 0,
		seed: 1,
		hash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	}, {
		path: "data_1",
		size: 256,
		seed: 1,
		hash: "a6452fbd8c12f8df622c1ca4c567f966801fb56442aca03b4e1303e7a412a9d5",
	}, {
		path: "sub/data_1",
		size: 256,
		seed: 1,
		hash: "a6452fbd8c12f8df622c1ca4c567f966801fb56442aca03b4e1303e7a412a9d5",
	}})
}
