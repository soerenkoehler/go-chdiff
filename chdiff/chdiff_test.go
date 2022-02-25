package chdiff

import (
	"crypto/sha256"
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
				AddNewEntry("a", "hashA", 1, time.Unix(1, 0))
	mockDigestCalculated = digest.
				NewDigest("dir_b", time.Now()).
				AddNewEntry("b", "hashB", 2, time.Unix(2, 0))
	mockDiffResult = diff.Diff{}
)

func TestRunSuite(t *testing.T) {
	testutil.RunSuite(t,
		func(t *testing.T) {
			mock = NewRegistry(t)
			mockDependencies = ChdiffDependencies{
				mockReader,
				mockWriter,
				mockCalculator,
				mockComparator,
				mockPrinter,
				mock.StdOut,
				mock.StdErr,
				func(e int) { mock.Register("exit", e) }}
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
					"out.txt",
					"SHA256")
			},
			"verify with path": func(t *testing.T) {
				testDigestVerify(t,
					[]string{"", "v", "x"},
					"x",
					"out.txt",
					"SHA256")
			},
			"create without path": func(t *testing.T) {
				testDigestCreate(t,
					[]string{"", "c"},
					".",
					"out.txt",
					"SHA256")
			},
			"create with path": func(t *testing.T) {
				testDigestCreate(t,
					[]string{"", "c", "x"},
					"x",
					"out.txt",
					"SHA256")
			}})
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
		Verify("calculate", Is(absDataPath), IsFunc(sha256.New)).
		Verify("compare", Is(mockDigestLoaded), Is(mockDigestCalculated)).
		Verify("print", Is(mockDiffResult)).
		NoMoreInvocations()
}

func testDigestCreate(
	t *testing.T,
	args []string,
	dataPath, digestPath, algorithm string) {

	absDataPath, _ := filepath.Abs(dataPath)
	absDigestPath := path.Join(absDataPath, digestPath)

	Chdiff("TEST", args, mockDependencies)

	mock.
		Verify("calculate", Is(absDataPath), IsFunc(sha256.New)).
		Verify("write", Is(absDigestPath), Is(mockDigestCalculated)).
		NoMoreInvocations()
}

func mockReader(path string) (digest.Digest, error) {
	mock.Register("read", path)
	return mockDigestLoaded, nil
}

func mockWriter(path string, digest digest.Digest) error {
	mock.Register("write", path, digest)
	return nil
}

func mockCalculator(path string, hashFactory digest.HashFactory) digest.Digest {
	mock.Register("calculate", path, hashFactory)
	return mockDigestCalculated
}

func mockComparator(old, new digest.Digest) diff.Diff {
	mock.Register("compare", old, new)
	return mockDiffResult
}

func mockPrinter(out io.Writer, diff diff.Diff) {
	mock.Register("print", diff)
}
