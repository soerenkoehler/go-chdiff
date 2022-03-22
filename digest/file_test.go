package digest_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/soerenkoehler/go-chdiff/digest"
)

func TestLoadNonexistantFile(t *testing.T) {
	_, err := digest.Load("../testdata/digest/file/path-without-digest", "algo")
	fmt.Println(err)
	if err.Error() != "open ../testdata/digest/file/path-without-digest/.chdiff.algo.txt: no such file or directory" {
		t.Fatal("expected: file not found error")
	}
}

func TestLoad(t *testing.T) {
	expected := digest.NewDigest("rootPath", "algo", time.Now())
	actual, err := digest.Load("../testdata/digest/file", "algo")
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(expected, actual) {
		t.Fatal(cmp.Diff(expected, actual))
	}
}

func TestSave(t *testing.T) {
	// TODO
	digest.NewDigest("rootPath", "algo", time.Now())
	t.Fatal("TODO")
}
