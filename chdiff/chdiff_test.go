package chdiff

import (
	"io"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/soerenkoehler/go-chdiff/diff"
	"github.com/soerenkoehler/go-chdiff/digest"
	. "github.com/soerenkoehler/go-testutils/mockutil"
	"github.com/soerenkoehler/go-testutils/testutil"
)

var (
	mock             Registry
	mockDependencies ChdiffDependencies

	mockDigestLoaded = digest.
				NewDigest("dir_a", time.Now()).
				AddNewEntry("a", "hashA")
	mockDigestCalculated = digest.
				NewDigest("dir_b", time.Now()).
				AddNewEntry("b", "hashB")
	mockDiffResult = diff.Diff{}
)

func TestRunSuite(t *testing.T) {
	testutil.RunSuite(t,
		func(t *testing.T) {
			mock = NewRegistry(t)
			mockDependencies = ChdiffDependencies{
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
		testutil.Suite{
			"no default": func(t *testing.T) {
				testErrorMessage(t,
					[]string{""},
					"error: expected one of \"create\",  \"verify\"\n")
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
					".chdiff.SHA256.txt",
					"SHA256")
			},
			"verify with path": func(t *testing.T) {
				testDigestVerify(t,
					[]string{"", "v", "x"},
					"x",
					".chdiff.SHA256.txt",
					"SHA256")
			},
			"create without path": func(t *testing.T) {
				testDigestCreate(t,
					[]string{"", "c"},
					".",
					".chdiff.SHA256.txt",
					"SHA256")
			},
			"create with path": func(t *testing.T) {
				testDigestCreate(t,
					[]string{"", "c", "x"},
					"x",
					".chdiff.SHA256.txt",
					"SHA256")
			},
		})
}

func testErrorMessage(
	t *testing.T,
	args []string,
	msg string) {

	Chdiff(
		"TEST",
		args,
		mockDependencies)
	if !strings.Contains(mock.StdErr.String(), msg) {
		t.Errorf("no or incorrect error message")
	}
}

func testDigestVerify(
	t *testing.T,
	args []string,
	dataPath, digestPath, algorithm string) {

	absDataPath, _ := filepath.Abs(dataPath)
	absDigestPath := path.Join(absDataPath, digestPath)

	Chdiff("TEST", args, mockDependencies)

	mock.
		Verify("read", Is(absDigestPath)).
		Verify("calculate", Is(absDataPath), Is(algorithm)).
		Verify("compare", Is(mockDigestLoaded), Is(mockDigestCalculated)).
		Verify("print", Is(mockDiffResult)).
		NoMoreInvocations()
}

func testDigestCreate(
	t *testing.T,
	args []string,
	dataPath, digestPath, algorithm string) {

	absDataPath, _ := filepath.Abs(dataPath)

	Chdiff("TEST", args, mockDependencies)

	mock.
		Verify("calculate", Is(absDataPath), Is(algorithm)).
		Verify("write", Is(mockDigestCalculated)).
		NoMoreInvocations()
}

func mockReader(path string) (digest.Digest, error) {
	mock.Register("read", path)
	return mockDigestLoaded, nil
}

func mockWriter(digest digest.Digest) error {
	mock.Register("write", digest)
	return nil
}

func mockCalculator(path, algorithm string) digest.Digest {
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
