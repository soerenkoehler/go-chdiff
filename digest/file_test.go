package digest_test

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/soerenkoehler/go-chdiff/digest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TSFile struct {
	suite.Suite
}

func TestSuiteFileRunner(t *testing.T) {
	suite.Run(t, &TSFile{})
}

func (s *TSFile) SetupTest() {
}

func (s *TSFile) TestLoadNonexistantFile() {
	path := "../testdata/digest/file"
	file := filepath.Join(path, "nonexistant-digest-file.txt")
	_, err := digest.Load(path, file)
	assert.EqualError(s.T(), err, fmt.Sprintf("lstat %v: no such file or directory", file))
}

func (s *TSFile) TestLoadBadDigest1Column() {
	path := "../testdata/digest/file"
	file := filepath.Join(path, "bad-digest-1-column.txt")
	_, err := digest.Load(path, file)
	assert.EqualError(s.T(), err, "invalid digest file")
}

func (s *TSFile) TestSaveLoad256() {
	s.testSaveLoad(32)
}

func (s *TSFile) TestSaveLoad512() {
	s.testSaveLoad(64)
}

func (s *TSFile) testSaveLoad(hashsize int) {
	digestPath := s.T().TempDir()
	digestTime := time.Now()
	digestFile := filepath.Join(digestPath, "test-digest.txt")
	expected := digest.NewDigest(digestPath, digestTime)
	expected.AddFileHash("file1", createRandomHash(hashsize))
	expected.AddFileHash("file2", createRandomHash(hashsize))
	digest.Save(expected, digestFile)
	actual, err := digest.Load(digestPath, digestFile)
	assert.Nil(s.T(), err)
	assert.True(s.T(), cmp.Equal(expected, actual), cmp.Diff(expected, actual))
}
