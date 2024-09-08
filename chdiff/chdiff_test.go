package chdiff_test

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/soerenkoehler/go-chdiff/chdiff"
	"github.com/soerenkoehler/go-chdiff/diff"
	"github.com/soerenkoehler/go-chdiff/digest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	mockDigestLoaded     = digest.Digest{Algorithm: digest.SHA256}
	mockDigestCalculated = digest.Digest{}
	mockDiffResult       = diff.Diff{}
)

type MockDependencies struct {
	mock.Mock
}

func (m *MockDependencies) DigestRead(dp, df string) (digest.Digest, error) {
	result := m.Called(dp, df)
	return result.Get(0).(digest.Digest), result.Error(1)
}

func (m *MockDependencies) DigestWrite(d digest.Digest, df string) error {
	return m.Called(d, df).Error(0)
}

func (m *MockDependencies) DigestCalculate(rp string, ht digest.HashType) digest.Digest {
	return m.Called(rp, ht).Get(0).(digest.Digest)
}

func (m *MockDependencies) DigestCompare(old, new digest.Digest) diff.Diff {
	return m.Called(old, new).Get(0).(diff.Diff)
}

func (m *MockDependencies) DiffPrint(out io.Writer, d diff.Diff) {
	m.Called(out, d)
}

func (m *MockDependencies) Stdout() io.Writer {
	return m.Called().Get(0).(io.Writer)
}

func (m *MockDependencies) Stderr() io.Writer {
	return m.Called().Get(0).(io.Writer)
}

func (m *MockDependencies) KongExit() func(int) {
	return m.Called().Get(0).(func(int))
}

type TSChdiff struct {
	suite.Suite
	Stdout       *strings.Builder
	Stderr       *strings.Builder
	Dependencies *MockDependencies
}

func TestSuiteRunner(t *testing.T) {
	suite.Run(t, &TSChdiff{})
}

func (s *TSChdiff) SetupTest() {
	s.Stdout = &strings.Builder{}
	s.Stderr = &strings.Builder{}
	s.Dependencies = &MockDependencies{}
	s.Dependencies.
		On("Stdout").Return(s.Stdout).Once().
		On("Stderr").Return(s.Stderr).Twice().
		On("KongExit").Return(
		func(e int) {
			s.Dependencies.MethodCalled("exit", e)
		})
}

func (s *TSChdiff) TestLoadConfig() {
	s.T().Setenv("HOME", "../testdata/chdiff/userhome")
	s.Dependencies.Mock.On("exit", mock.Anything).Return()

	chdiff.Chdiff("TEST", []string{""}, s.Dependencies)

	s.Dependencies.AssertExpectations(s.T())
	assert.Contains(
		s.T(), s.Stderr.String(),
		"[D] {Exclude:{Absolute:[] Relative:[] Anywhere:[]} LogLevel:debug}")
}

func (s *TSChdiff) TestLoadConfigCreateMissingFile() {
	tempUserhome := s.T().TempDir()
	s.T().Setenv("HOME", tempUserhome)
	s.Dependencies.Mock.On("exit", mock.Anything).Return()

	chdiff.Chdiff("TEST", []string{""}, s.Dependencies)

	s.Dependencies.AssertExpectations(s.T())
	assert.FileExists(s.T(), filepath.Join(tempUserhome, chdiff.UserConfigFileName))
	assert.Contains(
		s.T(), s.Stderr.String(),
		// Attention: Kong's error message contains double space between commands
		"error: expected one of \"create\",  \"verify\"\n")
}

func (s *TSChdiff) TestLoadConfigBadUserhome() {
	s.T().Setenv("HOME", "")
	s.Dependencies.Mock.On("exit", mock.Anything).Return()

	chdiff.Chdiff("TEST", []string{""}, s.Dependencies)

	s.Dependencies.AssertExpectations(s.T())
	assert.Contains(
		s.T(), s.Stderr.String(),
		"[W] can't determine user home")
}

func (s *TSChdiff) TestLoadConfigBadJson() {
	s.T().Setenv("HOME", "../testdata/chdiff/userhome-with-bad-config")
	s.Dependencies.Mock.On("exit", mock.Anything).Return()

	assert.PanicsWithError(s.T(),
		`/!\ reading config: invalid character 'i' looking for beginning of value`, func() {
			chdiff.Chdiff("TEST", []string{""}, s.Dependencies)
		})
}

func (s *TSChdiff) TestNoCommand() {
	testErrorMessage(s,
		[]string{""},
		// Attention: Kong's error message contains double space between commands
		"error: expected one of \"create\",  \"verify\"\n")
}

func (s *TSChdiff) TestUnknownCommand() {
	testErrorMessage(s,
		[]string{"", "bad-command"},
		"error: unexpected argument bad-command\n")
}

func (s *TSChdiff) TestVerifyWithoutPath() {
	testDigestVerify(s,
		[]string{"", "v"},
		".",
		chdiff.DefaultDigestName,
		digest.SHA256)
}

func (s *TSChdiff) TestVerifyWithPath() {
	testDigestVerify(s,
		[]string{"", "v", "x"},
		"x",
		chdiff.DefaultDigestName,
		digest.SHA256)
}

func (s *TSChdiff) TestDigestVerifyMissingDigestFile() {

	absDataPath, _ := filepath.Abs("x")
	absDigestFile := filepath.Join(absDataPath, chdiff.DefaultDigestName)

	s.Dependencies.
		On("exit", mock.Anything).Return().
		On("DigestRead", absDataPath, absDigestFile).Return(
		digest.Digest{},
		fmt.Errorf("no such file"))

	chdiff.Chdiff("TEST", []string{"", "v", "x"}, s.Dependencies)

	s.Dependencies.AssertExpectations(s.T())
	assert.Contains(s.T(), s.Stderr.String(), "[E] verify: no such file")
}

func (s *TSChdiff) TestDigestCreateSHA256DefaultName() {
	absDataPath, _ := filepath.Abs("x")
	absDigestFile := filepath.Join(absDataPath, chdiff.DefaultDigestName)

	s.Dependencies.
		On("DigestCalculate", absDataPath, digest.SHA256).Return(mockDigestCalculated).
		On("DigestWrite", mockDigestCalculated, absDigestFile).Return(nil)

	chdiff.Chdiff("TEST", []string{"", "c", "-a", "SHA256", "x"}, s.Dependencies)

	s.Dependencies.AssertExpectations(s.T())
}

func (s *TSChdiff) TestDigestCreateSHA256ExplicitName() {
	absDataPath, _ := filepath.Abs("x")
	absDigestPath, _ := filepath.Abs("y")
	absDigestFile := filepath.Join(absDigestPath, "explicit")

	s.Dependencies.
		On("DigestCalculate", absDataPath, digest.SHA256).Return(mockDigestCalculated).
		On("DigestWrite", mockDigestCalculated, absDigestFile).Return(nil)

	chdiff.Chdiff("TEST", []string{"", "c", "-a", "SHA256", "x", "-f", "y/explicit"}, s.Dependencies)

	s.Dependencies.AssertExpectations(s.T())
}

func (s *TSChdiff) TestDigestCreateSHA512() {
	absDataPath, _ := filepath.Abs("x")
	absDigestFile := filepath.Join(absDataPath, chdiff.DefaultDigestName)

	s.Dependencies.
		On("DigestCalculate", absDataPath, digest.SHA512).Return(mockDigestCalculated).
		On("DigestWrite", mockDigestCalculated, absDigestFile).Return(nil)

	chdiff.Chdiff("TEST", []string{"", "c", "-a", "SHA512", "x"}, s.Dependencies)

	s.Dependencies.AssertExpectations(s.T())
}

func (s *TSChdiff) TestDigestCreateBadAlgorithm() {
	s.Dependencies.Mock.On("exit", mock.Anything).Return()

	chdiff.Chdiff("TEST", []string{"", "c", "-a", "WRONG", "x"}, s.Dependencies)

	s.Dependencies.AssertExpectations(s.T())
	assert.Contains(s.T(), s.Stderr.String(), "--algorithm must be one of \"SHA256\",\"SHA512\" but got \"WRONG\"")
}

func testErrorMessage(
	s *TSChdiff,
	args []string,
	expected string) {

	s.Dependencies.Mock.On("exit", mock.Anything).Return()

	chdiff.Chdiff("TEST", args, s.Dependencies)

	s.Dependencies.Mock.AssertExpectations(s.T())
	assert.Contains(s.T(), s.Stderr.String(), expected)
}

func testDigestVerify(
	s *TSChdiff,
	args []string,
	dataPath, digestPath string,
	algorithm digest.HashType) {

	absDataPath, _ := filepath.Abs(dataPath)
	absDigestFile := filepath.Join(absDataPath, chdiff.DefaultDigestName)

	s.Dependencies.
		On("Stdout").Return(s.Stdout).Once().
		On("DigestRead", absDataPath, absDigestFile).Return(mockDigestLoaded, nil).
		On("DigestCalculate", absDataPath, digest.SHA256).Return(mockDigestCalculated).
		On("DigestCompare", mockDigestLoaded, mockDigestCalculated).Return(mockDiffResult).
		On("DiffPrint", s.Stdout, mockDiffResult).Return()

	chdiff.Chdiff("TEST", args, s.Dependencies)

	s.Dependencies.AssertExpectations(s.T())
}
