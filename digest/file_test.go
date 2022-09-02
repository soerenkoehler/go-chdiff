package digest_test

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/soerenkoehler/go-chdiff/digest"
)

func TestLoadNonexistantFile(t *testing.T) {
	path := "../testdata/digest/file/path-without-digest"
	file := filepath.Join(path, ".chdiff.txt")
	expectedError := fmt.Sprintf("lstat %v: no such file or directory", file)
	_, err := digest.Load(path, file)
	if err == nil || err.Error() != expectedError {
		t.Fatalf("\nexpected: %v\n  actual: %v", expectedError, err)
	}
}

func TestSaveLoad256(t *testing.T) {
	testSaveLoad(t, "SHA256")
}

func TestSaveLoad512(t *testing.T) {
	testSaveLoad(t, "SHA512")
}

func testSaveLoad(t *testing.T, algorithm string) {
	digestPath := t.TempDir()
	digestTime := time.Now()
	digestFile := digest.DefaultDigestFile(digestPath)
	expected := digest.NewDigest(digestPath, digestTime)
	digest.Save(expected, digestFile)
	actual, err := digest.Load(digestPath, digestFile)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(expected, actual) {
		t.Fatal(cmp.Diff(expected, actual))
	}
}
