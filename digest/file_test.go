package digest_test

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/soerenkoehler/go-chdiff/digest"
	"github.com/stretchr/testify/assert"
)

func TestLoadNonexistantFile(t *testing.T) {
	path := "../testdata/digest/file/path-without-digest"
	file := filepath.Join(path, "test-digest.txt")
	_, err := digest.Load(path, file)
	assert.EqualError(t, err, fmt.Sprintf("lstat %v: no such file or directory", file))
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
	expected.AddFileHash("file1", createRandomHash(hashsize))
	expected.AddFileHash("file2", createRandomHash(hashsize))
	digest.Save(expected, digestFile)
	actual, err := digest.Load(digestPath, digestFile)
	assert.Nil(t, err)
	assert.True(t, cmp.Equal(expected, actual), cmp.Diff(expected, actual))
}
