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

func TestWrongAlgorithm(t *testing.T) {
	if len(createDigest(t, []testCase{{
		path: "invalid",
		size: 0,
		seed: 1,
		hash: "invalid",
	}}, "INVALID")) != 0 {
		t.Fatal("invalid algorithm must not create digest entries")
	}
}

func TestDigest256(t *testing.T) {
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
	}}, "SHA256")
}

func TestDigest512(t *testing.T) {
	runDigestCalculationTest(t, []testCase{{
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

func runDigestCalculationTest(
	t *testing.T,
	data []testCase,
	algorithm string) {

	verifyDigest(t, data,
		createDigest(t, data, algorithm))
}

func createDigest(
	t *testing.T,
	data []testCase,
	algorithm string) Digest {

	root := t.TempDir()

	for _, dataPoint := range data {
		file := path.Join(root, dataPoint.path)
		datautil.CreateRandomFile(file, dataPoint.size, dataPoint.seed)
	}

	return calculateDigest(root, algorithm)
}

func verifyDigest(
	t *testing.T,
	data []testCase,
	digest Digest) {

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
			t.Errorf("actual hash (%v) does not match expected hash (%v) (test file: %v)",
				actualHash,
				expectedHash,
				expectedPath)
		}
	}
}
