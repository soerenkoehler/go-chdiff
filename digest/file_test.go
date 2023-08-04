package digest_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/soerenkoehler/go-chdiff/digest"
)

func TestLoadNonexistantFile(t *testing.T) {
	path := "../testdata/digest/file/path-without-digest"
	file := filepath.Join(path, "test-digest.txt")
	expectedError := fmt.Sprintf("lstat %v: no such file or directory", file)
	_, err := digest.Load(path, file)
	if err == nil || err.Error() != expectedError {
		t.Fatalf("\nexpected: %v\n  actual: %v", expectedError, err)
	}
}

func TestSaveLoad256(t *testing.T) {
	testSaveLoad(t, 32)
}

func TestSaveLoad512(t *testing.T) {
	testSaveLoad(t, 64)
}

func testSaveLoad(t *testing.T, hashsize int) {
	digestPath := t.TempDir()
	digestTime := time.Now()
	digestFile := filepath.Join(digestPath, "test-digest.txt")
	expected := digest.NewDigest(digestPath, digestTime)
	expected.AddFileHash("file", createRandomHash(hashsize))
	digest.Save(expected, digestFile)
	actual, err := digest.Load(digestPath, digestFile)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(expected, actual) {
		t.Fatal(cmp.Diff(expected, actual))
	}
}

func createRandomHash(hashsize int) string {
	hashbytes := make([]byte, hashsize)
	rand.Read(hashbytes)
	return hex.EncodeToString(hashbytes)
}
