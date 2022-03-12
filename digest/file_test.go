package digest

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestLoadNonexistantFile(t *testing.T) {
	_, err := Load("../testdata/digest/file/non-existant-file")
	if err == nil {
		t.Fatal("expected: file not found error")
	}
}

func TestLoad(t *testing.T) {
	expected := newDigest("rootPath", "algo", time.Now())
	actual, err := Load("../testdata/digest/file/digest.txt")
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(expected, actual) {
		t.Fatal(cmp.Diff(expected, actual))
	}
}

func TestSave(t *testing.T) {
	// TODO
	newDigest("rootPath", "algo", time.Now())
	t.Fatal("TODO")
}
