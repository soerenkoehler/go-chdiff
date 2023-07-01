package chdiff_test

import (
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/soerenkoehler/go-chdiff/chdiff"
	"github.com/soerenkoehler/go-chdiff/diff"
	"github.com/soerenkoehler/go-chdiff/digest"
	. "github.com/soerenkoehler/go-util-test/mock"
	"github.com/soerenkoehler/go-util-test/test"
)

var (
	mock             Registry
	mockDependencies chdiff.ChdiffDependencies

	mockDigestLoaded = digest.Digest{
		Algorithm: digest.SHA256}
	mockDigestCalculated = digest.Digest{}
	mockDiffResult       = diff.Diff{}
)

func TestRunSuite(t *testing.T) {
	test.RunSuite(t,
		func(t *testing.T) {
			mock = NewRegistry(t)
			mockDependencies = chdiff.ChdiffDependencies{
				DigestRead:      mockReader,
				DigestWrite:     mockWriter,
				DigestCalculate: mockCalculator,
				DigestCompare:   mockComparator,
				DiffPrint:       mockPrinter,
				Stdout:          mock.StdOut,
				Stderr:          mock.StdErr,
				KongExit:        func(e int) { mock.Register("exit", e) }}
		},
		nil,
		test.Suite{
			"no default": func(t *testing.T) {
				testErrorMessage(t,
					[]string{""},
					"error: expected one of \"create\",  \"verify\"\n")
				// Attention: Kong's error message contains double space
				// between commands
			},
			"unknown command": func(t *testing.T) {
				testErrorMessage(t,
					[]string{"", "bad-command"},
					"error: unexpected argument bad-command\n")
			},
			"verify without path": func(t *testing.T) {
				testDigestVerify(t,
					[]string{"", "v"},
					".",
					chdiff.DefaultDigestName,
					digest.SHA256)
			},
			"verify with path": func(t *testing.T) {
				testDigestVerify(t,
					[]string{"", "v", "x"},
					"x",
					chdiff.DefaultDigestName,
					digest.SHA256)
			},
			"create without path": func(t *testing.T) {
				testDigestCreate(t,
					[]string{"", "c"},
					".",
					chdiff.DefaultDigestName,
					digest.SHA256)
			},
			"create with path": func(t *testing.T) {
				testDigestCreate(t,
					[]string{"", "c", "x"},
					"x",
					chdiff.DefaultDigestName,
					digest.SHA256)
			},
		})
}

func testErrorMessage(
	t *testing.T,
	args []string,
	expected string) {

	chdiff.Chdiff(
		"TEST",
		args,
		mockDependencies)
	actual := mock.StdErr.String()
	if !strings.Contains(actual, expected) {
		t.Errorf("\nexpected: %v\nactual: %v", expected, actual)
	}
}

func testDigestVerify(
	t *testing.T,
	args []string,
	dataPath, digestPath string,
	algorithm digest.HashType) {

	absDataPath, _ := filepath.Abs(dataPath)
	absDigestFile := filepath.Join(absDataPath, chdiff.DefaultDigestName)

	chdiff.Chdiff("TEST", args, mockDependencies)

	mock.
		Verify("read", Is(absDataPath), Is(absDigestFile)).
		Verify("calculate", Is(absDataPath), Is(algorithm)).
		Verify("compare", Is(mockDigestLoaded), Is(mockDigestCalculated)).
		Verify("print", Is(mockDiffResult)).
		NoMoreInvocations()
}

func testDigestCreate(
	t *testing.T,
	args []string,
	dataPath, digestPath string,
	algorithm digest.HashType) {

	absDataPath, _ := filepath.Abs(dataPath)
	absDigestFile := filepath.Join(absDataPath, chdiff.DefaultDigestName)

	chdiff.Chdiff("TEST", args, mockDependencies)

	mock.
		Verify("calculate", Is(absDataPath), Is(algorithm)).
		Verify("write", Is(mockDigestCalculated), Is(absDigestFile)).
		NoMoreInvocations()
}

func mockReader(path, file string) (digest.Digest, error) {
	mock.Register("read", path, file)
	return mockDigestLoaded, nil
}

func mockWriter(digest digest.Digest, digestFile string) error {
	mock.Register("write", digest, digestFile)
	return nil
}

func mockCalculator(path string, algorithm digest.HashType) digest.Digest {
	mock.Register("calculate", path, algorithm)
	return mockDigestCalculated
}

func mockComparator(old, new digest.Digest) diff.Diff {
	mock.Register("compare", old, new)
	return mockDiffResult
}

func mockPrinter(out io.Writer, diff diff.Diff) {
	mock.Register("print", diff)
}
