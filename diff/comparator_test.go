package diff_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/soerenkoehler/go-chdiff/common"
	"github.com/soerenkoehler/go-chdiff/diff"
	"github.com/soerenkoehler/go-chdiff/digest"
	"github.com/stretchr/testify/suite"
)

const (
	digestPath1 = "/path/to/digestfile"
	digestPath2 = "/path/to/dir"
	digestFile1 = "../testdata/diff/comparator/digest-old.txt"
	digestFile2 = "../testdata/diff/comparator/digest-new.txt"
	fileHash1   = "hash1"
	fileHash2   = "hash2"
)

var (
	digestTime1 time.Time = time.Date(2020, 3, 4, 16, 17, 18, 0, time.Local)
	digestTime2 time.Time = time.Date(2022, 1, 2, 13, 14, 15, 0, time.Local)
)

type TSComparator struct {
	suite.Suite
	Stdout *strings.Builder
}

func TestSuiteRunner(t *testing.T) {
	suite.Run(t, &TSComparator{})
}

func (s *TSComparator) SetupTest() {
	s.Stdout = &strings.Builder{}
}

func (s *TSComparator) TestOutputEmptyDiff() {
	diff.Print(s.Stdout, makeDiff(0, 0, 0, 0), false)

	expect(s, []string{}, 0, 0, 0, 0)
}

func (s *TSComparator) TestOutputNoChanges() {
	diff.Print(s.Stdout, makeDiff(2, 0, 0, 0), false)

	expect(s, []string{}, 2, 0, 0, 0)
}

func (s *TSComparator) TestOutputNoChangesShowIdentical() {
	diff.Print(s.Stdout, makeDiff(2, 0, 0, 0), true)

	expect(s, []string{
		"= relPath0",
		"= relPath1",
	}, 2, 0, 0, 0)
}

func (s *TSComparator) TestOutputWithChanges() {
	diff.Print(s.Stdout, makeDiff(4, 3, 5, 7), false)

	expect(s, []string{
		"+ relPath10",
		"+ relPath11",
		"- relPath12",
		"- relPath13",
		"- relPath14",
		"- relPath15",
		"- relPath16",
		"- relPath17",
		"- relPath18",
		"M relPath4",
		"M relPath5",
		"M relPath6",
		"+ relPath7",
		"+ relPath8",
		"+ relPath9",
	}, 4, 3, 5, 7)
}

func (s *TSComparator) TestOutputWithChangesShowIdentical() {
	diff.Print(s.Stdout, makeDiff(4, 3, 5, 7), true)

	expect(s, []string{
		"= relPath0",
		"= relPath1",
		"+ relPath10",
		"+ relPath11",
		"- relPath12",
		"- relPath13",
		"- relPath14",
		"- relPath15",
		"- relPath16",
		"- relPath17",
		"- relPath18",
		"= relPath2",
		"= relPath3",
		"M relPath4",
		"M relPath5",
		"M relPath6",
		"+ relPath7",
		"+ relPath8",
		"+ relPath9",
	}, 4, 3, 5, 7)
}

func (s *TSComparator) TestCompare() {
	diff.Print(s.Stdout, diff.Compare(
		makeDigest(s, digestPath1, digestFile1, digestTime1),
		makeDigest(s, digestPath2, digestFile2, digestTime2)),
		false)

	expect(s,
		[]string{
			"- f0",
			"M f2",
			"+ f3",
		}, 1, 1, 1, 1)
}

func makeDiff(identical, modified, added, removed int32) diff.Diff {
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

func expect(s *TSComparator, entries []string, identical, modified, added, removed int32) {
	// for non-empty entries list require a final newline
	entriesText := strings.Join(append(entries, ""), "\n")

	expected := fmt.Sprintf(
		"Old: (%s) %v\nNew: (%s) %v\n%vIdentical: %v | Modified: %v | Added: %v | Removed: %v\n",
		common.LocationTimeFormat.FormatString(digestTime1), digestPath1,
		common.LocationTimeFormat.FormatString(digestTime2), digestPath2,
		entriesText,
		identical, modified, added, removed)

	actual := s.Stdout.String()

	if actual != expected {
		s.T().Fatalf("expected:\n%v\nactual:\n%v", expected, actual)
	}
}

func makeDigest(s *TSComparator, digestPath, digestFile string, modTime time.Time) digest.Digest {
	os.Chtimes(digestFile, modTime, modTime)
	result, err := digest.Load(digestPath, digestFile)
	if err != nil {
		s.T().Fatal(err)
	}
	return result
}
