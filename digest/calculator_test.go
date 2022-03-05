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

func TestDigest256(t *testing.T) {
	verifyDigest(t, []testCase{{
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
	}}, "SHA256")
}

func TestDigest512(t *testing.T) {
	verifyDigest(t, []testCase{{
		path: "zero",
		size: 0,
		seed: 1,
		hash: "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
	}, {
		path: "data_1",
		size: 256,
		seed: 1,
		hash: "f3f00e46e5dc3819b8268afedb1221f25a4c29d3223979ede1df107155cc75bd427a5795b820fbd83fd4785899cb9de201b770a2c88a3bed90be37e82156e10b",
	}, {
		path: "sub/data_1",
		size: 256,
		seed: 1,
		hash: "f3f00e46e5dc3819b8268afedb1221f25a4c29d3223979ede1df107155cc75bd427a5795b820fbd83fd4785899cb9de201b770a2c88a3bed90be37e82156e10b",
	}}, "SHA512")
}

func verifyDigest(
	t *testing.T,
	data []testCase,
	algorithm string) {

	digest := Calculate(createData(t, data), algorithm)

	if len(*digest.Entries) != len(data) {
		t.Fatal("Digest size must match number of input data points")
	}

	for _, dataPoint := range data {
		verifyDataPoint(t, dataPoint, (*digest.Entries)[dataPoint.path])
	}
}

func createData(
	t *testing.T,
	data []testCase) string {

	root := t.TempDir()

	for _, dataPoint := range data {
		file := path.Join(root, dataPoint.path)
		datautil.CreateRandomFile(file, dataPoint.size, dataPoint.seed)
	}

	return root
}

func verifyDataPoint(
	t *testing.T,
	dataPoint testCase,
	digestEntry digestEntry) {

	expectedPath := dataPoint.path
	actualPath := digestEntry.file
	if actualPath != expectedPath {
		t.Errorf("DigestEntry.file (%v) must match Digest map key (%v)",
			actualPath,
			expectedPath)
	}

	expectedHash := dataPoint.hash
	actualHash := digestEntry.Hash
	if actualHash != expectedHash {
		t.Errorf("hash mismatch\nexpected: %v\nactual: %v\ntest file: %v",
			expectedHash,
			actualHash,
			expectedPath)
	}
}
