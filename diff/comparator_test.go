package diff_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/soerenkoehler/go-chdiff/common"
	"github.com/soerenkoehler/go-chdiff/diff"
	"github.com/soerenkoehler/go-chdiff/digest"
	. "github.com/soerenkoehler/go-util-test/mock"
	"github.com/soerenkoehler/go-util-test/test"
)

const (
	digestPath1 = "/path/to/digestfile"
	digestPath2 = "/path/to/dir"
	fileHash1   = "hash1"
	fileHash2   = "hash2"
)

var (
	digestTime1 time.Time = time.Date(2020, 3, 4, 16, 17, 18, 0, time.Local)
	digestTime2 time.Time = time.Date(2022, 1, 2, 13, 14, 15, 0, time.Local)
	mock        Registry
)

func TestRunSuite(t *testing.T) {
	test.RunSuite(t,
		func(t *testing.T) {
			mock = NewRegistry(t)
		},
		nil,
		test.Suite{
			"print empty diff": func(t *testing.T) {
				diff.Print(mock.StdOut, makeDiff(t, 0, 0, 0, 0))

				expect(t, []string{}, 0, 0, 0, 0)
			},

			"print diff with no changes": func(t *testing.T) {
				diff.Print(mock.StdOut, makeDiff(t, 2, 0, 0, 0))

				expect(t, []string{}, 2, 0, 0, 0)
			},

			"print diff with changes": func(t *testing.T) {
				diff.Print(mock.StdOut, makeDiff(t, 0, 3, 5, 7))

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
				diff.Print(mock.StdOut, diff.Compare(makeDigests(t)))

				expect(t,
					[]string{
						"- f0",
						"* f2",
						"+ f3",
					}, 1, 1, 1, 1)
			},
		})
}

// TODO unused?
// func parseTime(t *testing.T, s string) time.Time {
// 	time, err := time.Parse(common.LocationTimeFormat, s)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	return time
// }

func makeDiff(t *testing.T, identical, modified, added, removed int32) diff.Diff {
	result := diff.Diff{
		LocationA: common.Location{
			Path: digestPath1,
			Time: digestTime1},
		LocationB: common.Location{
			Path: digestPath2,
			Time: digestTime2},
		Entries: map[string]diff.DiffEntry{}}
	entry := 0
	add := func(count int32, status diff.DiffStatus) {
		for ; count > 0; count-- {
			relPath := fmt.Sprintf("relPath%d", entry)
			result.Entries[relPath] = diff.DiffEntry{
				File:   relPath,
				Status: status,
			}
			entry++
		}
	}
	add(identical, diff.Identical)
	add(modified, diff.Modified)
	add(added, diff.Added)
	add(removed, diff.Removed)
	return result
}

func expect(t *testing.T, entries []string, identical, modified, added, removed int32) {
	// for non-empty entries list require a final newline
	entriesText := strings.Join(append(entries, ""), "\n")

	expected := fmt.Sprintf(
		"Old: (%s) %v\nNew: (%s) %v\n%vIdentical: %v | Modified: %v | Added: %v | Removed: %v\n",
		common.LocationTimeFormat.FormatString(digestTime1), digestPath1,
		common.LocationTimeFormat.FormatString(digestTime2), digestPath2,
		entriesText,
		identical, modified, added, removed)

	actual := mock.StdOut.String()

	if actual != expected {
		t.Fatalf("expected:\n%v\nactual:\n%v", expected, actual)
	}
}

func makeDigests(t *testing.T) (digest.Digest, digest.Digest) {
	d1, err := digest.Load(
		"../testdata/diff/comparator/",
		"../testdata/diff/comparator/digest-old.txt")
	if err != nil {
		t.Fatal(err)
	}
	d2, err := digest.Load(
		"../testdata/diff/comparator/",
		"../testdata/diff/comparator/digest-new.txt")
	if err != nil {
		t.Fatal(err)
	}
	return d1, d2
}
