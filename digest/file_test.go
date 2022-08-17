package digest_test

import (
	"fmt"
	"path"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/soerenkoehler/go-chdiff/digest"
)

func TestLoadNonexistantFile(t *testing.T) {
	_, err := digest.Load("../testdata/digest/file/path-without-digest", "algo")
	fmt.Println(err)
	if err == nil || err.Error() != "open ../testdata/digest/file/path-without-digest/.chdiff.algo.txt: no such file or directory" {
		t.Fatal("expected: file not found error")
	}
}

func TestSaveLoad(t *testing.T) {
	path := path.Join(t.TempDir(), "rootPath")
	expected := digest.NewDigest(path, "algo", time.Now())
	digest.Save(expected)
	actual, err := digest.Load(path, "algo")
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(expected, actual) {
		t.Fatal(cmp.Diff(expected, actual))
	}
}
