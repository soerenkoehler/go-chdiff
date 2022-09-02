package digest_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/soerenkoehler/go-chdiff/digest"
	"github.com/soerenkoehler/go-chdiff/util"
)

func TestLoadNonexistantFile(t *testing.T) {
	expectedError := "lstat ../testdata/digest/file/path-without-digest/.chdiff.txt: no such file or directory"
	_, err := digest.Load("../testdata/digest/file/path-without-digest")
	if err == nil || err.Error() != expectedError {
		t.Fatalf("\nexpected: %v\n  actual: %v", expectedError, err)
	}
}

func TestSaveLoad(t *testing.T) {
	digestPath := t.TempDir()
	digestTime := time.Now()
	digestFile:=
	expected := digest.NewDigest(digestPath, "algo", digestTime)
	digest.Save(expected)
	actual, err := digest.Load(digestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(expected, actual) {
		t.Fatal(cmp.Diff(expected, actual))
	}
}
