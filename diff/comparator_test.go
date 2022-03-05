package diff

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/soerenkoehler/go-chdiff/common"
	"github.com/soerenkoehler/go-chdiff/digest"
	. "github.com/soerenkoehler/go-testutils/mockutil"
	"github.com/soerenkoehler/go-testutils/testutil"
)

const (
	rootPath1    = "/path/to/digestfile"
	rootPath2    = "/path/to/dir"
	rootTimeStr1 = "2020-03-04 16:17:18"
	rootTimeStr2 = "2022-01-02 13:14:15"
	fileHash1    = "hash1"
	fileHash2    = "hash2"
)

var mock Registry

func TestRunSuite(t *testing.T) {
	testutil.RunSuite(t,
		func(t *testing.T) {
			mock = NewRegistry(t)
		},
		nil,
		testutil.Suite{
			"print empty diff": func(t *testing.T) {
				Print(mock.StdOut, makeDiff(t, 0, 0, 0, 0))

				expect(t, []string{}, 0, 0, 0, 0)
			},

			"print diff with no changes": func(t *testing.T) {
				Print(mock.StdOut, makeDiff(t, 2, 0, 0, 0))

				expect(t, []string{}, 2, 0, 0, 0)
			},

			"print diff with changes": func(t *testing.T) {
				Print(mock.StdOut, makeDiff(t, 0, 3, 5, 7))

				expect(t, []string{
					"* relPath0",
					"* relPath1",
					"- relPath10",
					"- relPath11",
					"- relPath12",
					"- relPath13",
					"- relPath14",
					"* relPath2",
					"+ relPath3",
					"+ relPath4",
					"+ relPath5",
					"+ relPath6",
					"+ relPath7",
					"- relPath8",
					"- relPath9",
				}, 0, 3, 5, 7)
			},
			"compare by hash": func(t *testing.T) {
				Print(mock.StdOut, Compare(makeDigests(t)))

				expect(t,
					[]string{
						"- f0",
						"* f5",
						"* f6",
						"* f7",
						"* f8",
						"+ f9",
					}, 4, 4, 1, 1)
			},
		})
}

func parseTime(t *testing.T, s string) time.Time {
	time, err := time.Parse(common.LocationTimeFormat, s)
	if err != nil {
		t.Fatal(err)
	}
	return time
}

func makeDiff(t *testing.T, identical, modified, added, removed int32) Diff {
	result := Diff{
		locationA: common.Location{
			Path: rootPath1,
			Time: parseTime(t, rootTimeStr1)},
		locationB: common.Location{
			Path: rootPath2,
			Time: parseTime(t, rootTimeStr2)},
		entries: diffEntries{}}
	entry := 0
	add := func(count int32, status diffStatus) {
		for ; count > 0; count-- {
			relPath := fmt.Sprintf("relPath%d", entry)
			result.entries[relPath] = diffEntry{
				file:   relPath,
				status: status,
			}
			entry++
		}
	}
	add(identical, Identical)
	add(modified, Modified)
	add(added, Added)
	add(removed, Removed)
	return result
}

func expect(t *testing.T, entries []string, identical, modified, added, removed int32) {
	// for non-empty entries list require a final newline
	entriesText := strings.Join(append(entries, ""), "\n")

	expected := fmt.Sprintf(
		"Old: (%v) %v\nNew: (%v) %v\n%vIdentical: %v | Modified: %v | Added: %v | Removed: %v\n",
		rootTimeStr1, rootPath1, rootTimeStr2, rootPath2,
		entriesText,
		identical, modified, added, removed)

	actual := mock.StdOut.String()

	if actual != expected {
		t.Fatalf("should output:\n%v\nbut got:\n%v", expected, actual)
	}
}

func makeDigests(t *testing.T) (digest.Digest, digest.Digest) {
	d1, err := digest.Load("../testdata/diff/comparator/digest-old.txt")
	if err != nil {
		t.Fatal(err)
	}
	d2, err := digest.Load("../testdata/diff/comparator/digest-new.txt")
	if err != nil {
		t.Fatal(err)
	}
	return d1, d2
}
