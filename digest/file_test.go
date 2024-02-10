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

type TestSuiteFile struct {
	suite.Suite
}

func TestSuiteFileRunner(t *testing.T) {
	suite.Run(t, &TestSuiteFile{})
}

func (s *TestSuiteFile) SetupTest() {
}

func (s *TestSuiteFile) TestLoadNonexistantFile() {
	path := "../testdata/digest/file"
	file := filepath.Join(path, "nonexistant-digest-file.txt")
	_, err := digest.Load(path, file)
	assert.EqualError(s.T(), err, fmt.Sprintf("lstat %v: no such file or directory", file))
}

func (s *TestSuiteFile) TestLoadBadDigest1Column() {
	path := "../testdata/digest/file"
	file := filepath.Join(path, "bad-digest-1-column.txt")
	_, err := digest.Load(path, file)
	assert.EqualError(s.T(), err, "invalid digest file")
}

func (s *TestSuiteFile) TestSaveLoad256() {
	s.testSaveLoad(32)
}

func (s *TestSuiteFile) TestSaveLoad512() {
	s.testSaveLoad(64)
}

func (s *TestSuiteFile) testSaveLoad(hashsize int) {
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
