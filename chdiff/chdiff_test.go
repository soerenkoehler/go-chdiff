package chdiff_test

import (
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

type TestSuite struct {
	suite.Suite
	Stdout       *strings.Builder
	Stderr       *strings.Builder
	Dependencies *MockDependencies
}

func (s *TestSuite) SetupTest() {
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

func TestSuiteRunner(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) TestNoCommand() {
	testErrorMessage(s,
		[]string{""},
		// Attention: Kong's error message contains double space between commands
		"error: expected one of \"create\",  \"verify\"\n")
}

func (s *TestSuite) TestUnknownCommand() {
	testErrorMessage(s,
		[]string{"", "bad-command"},
		"error: unexpected argument bad-command\n")
}

func (s *TestSuite) TestVerifyWithoutPath() {
	testDigestVerify(s,
		[]string{"", "v"},
		".",
		chdiff.DefaultDigestName,
		digest.SHA256)
}

func (s *TestSuite) TestVerifyWithPath() {
	testDigestVerify(s,
		[]string{"", "v", "x"},
		"x",
		chdiff.DefaultDigestName,
		digest.SHA256)
}

func (s *TestSuite) TestDigestCreateSHA256DefaultName() {
	absDataPath, _ := filepath.Abs("x")
	absDigestFile := filepath.Join(absDataPath, chdiff.DefaultDigestName)

	s.Dependencies.
		On("DigestCalculate", absDataPath, digest.SHA256).Return(mockDigestCalculated).
		On("DigestWrite", mockDigestCalculated, absDigestFile).Return(nil)

	chdiff.Chdiff("TEST", []string{"", "c", "-a", "SHA256", "x"}, s.Dependencies)

	s.Dependencies.AssertExpectations(s.T())
}

func (s *TestSuite) TestDigestCreateSHA256ExplicitName() {
	absDataPath, _ := filepath.Abs("x")
	absDigestPath, _ := filepath.Abs("y")
	absDigestFile := filepath.Join(absDigestPath, "explicit")

	s.Dependencies.
		On("DigestCalculate", absDataPath, digest.SHA256).Return(mockDigestCalculated).
		On("DigestWrite", mockDigestCalculated, absDigestFile).Return(nil)

	chdiff.Chdiff("TEST", []string{"", "c", "-a", "SHA256", "x", "-f", "y/explicit"}, s.Dependencies)

	s.Dependencies.AssertExpectations(s.T())
}

func (s *TestSuite) TestDigestCreateSHA512() {
	absDataPath, _ := filepath.Abs("x")
	absDigestFile := filepath.Join(absDataPath, chdiff.DefaultDigestName)

	s.Dependencies.
		On("DigestCalculate", absDataPath, digest.SHA512).Return(mockDigestCalculated).
		On("DigestWrite", mockDigestCalculated, absDigestFile).Return(nil)

	chdiff.Chdiff("TEST", []string{"", "c", "-a", "SHA512", "x"}, s.Dependencies)

	s.Dependencies.AssertExpectations(s.T())
}

func (s *TestSuite) TestDigestCreateBadAlgorithm() {
	s.Dependencies.Mock.On("exit", mock.Anything).Return()

	chdiff.Chdiff("TEST", []string{"", "c", "-a", "WRONG", "x"}, s.Dependencies)

	s.Dependencies.AssertExpectations(s.T())
	assert.Contains(s.T(), s.Stderr.String(), "--algorithm must be one of \"SHA256\",\"SHA512\" but got \"WRONG\"")
}

func testErrorMessage(
	s *TestSuite,
	args []string,
	expected string) {

	s.Dependencies.Mock.On("exit", mock.Anything).Return()

	chdiff.Chdiff(
		"TEST",
		args,
		s.Dependencies)

	s.Dependencies.Mock.AssertExpectations(s.T())
	assert.Contains(s.T(), s.Stderr.String(), expected)
}

func testDigestVerify(
	s *TestSuite,
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

// t *testing.T,
// args []string,
// dataPath, digestPath string,
// algorithm digest.HashType) {

// absDataPath, _ := filepath.Abs(dataPath)
// absDigestFile := filepath.Join(absDataPath, chdiff.DefaultDigestName)

// chdiff.Chdiff("TEST", args, mockDependencies)

// mock.
// 	Verify("calculate", Is(absDataPath), Is(algorithm)).
// 	Verify("write", Is(mockDigestCalculated), Is(absDigestFile)).
// 	NoMoreInvocations()

// func mockReader(path, file string) (digest.Digest, error) {
// 	mock.Register("read", path, file)
// 	return mockDigestLoaded, nil
// }

// func mockWriter(digest digest.Digest, digestFile string) error {
// 	mock.Register("write", digest, digestFile)
// 	return nil
// }

// func mockCalculator(path string, algorithm digest.HashType) digest.Digest {
// 	mock.Register("calculate", path, algorithm)
// 	return mockDigestCalculated
// }

// func mockComparator(old, new digest.Digest) diff.Diff {
// 	mock.Register("compare", old, new)
// 	return mockDiffResult
// }

// func mockPrinter(out io.Writer, diff diff.Diff) {
// 	mock.Register("print", diff)
// }
