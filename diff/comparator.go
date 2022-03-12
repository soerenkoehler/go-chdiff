package diff

import (
	"fmt"
	"io"
	"sort"

	"github.com/soerenkoehler/go-chdiff/common"
	"github.com/soerenkoehler/go-chdiff/digest"
)

type Comparator func(digest.Digest, digest.Digest) Diff

type DiffPrinter func(io.Writer, Diff)

var statusIcon map[DiffStatus]string = map[DiffStatus]string{
	Identical: " ",
	Modified:  "*",
	Added:     "+",
	Removed:   "-"}

func Compare(old, new digest.Digest) Diff {
	diffEntries := DiffEntries{}

	// step 1: identical, modified and removed files
	for path, oldEntry := range *old.Entries {
		status := Removed
		if newEntry, newExists := (*new.Entries)[path]; newExists {
			if oldEntry.Hash == newEntry.Hash {
				status = Identical
			} else {
				status = Modified
			}
		}
		diffEntries[path] = DiffEntry{
			File:   path,
			Status: status}
	}

	// step 2: added files
	for path := range *new.Entries {
		if _, oldExists := (*old.Entries)[path]; !oldExists {
			diffEntries[path] = DiffEntry{
				File:   path,
				Status: Added}
		}
	}

	return Diff{
		LocationA: old.Location,
		LocationB: new.Location,
		Entries:   diffEntries}
}

func Print(out io.Writer, diff Diff) {
	fmt.Fprintf(out,
		"Old: (%s) %v\nNew: (%s) %v\n",
		diff.LocationA.Time.Format(common.LocationTimeFormat),
		diff.LocationA.Path,
		diff.LocationB.Time.Format(common.LocationTimeFormat),
		diff.LocationB.Path)

	count := make(map[DiffStatus]int32, 4)

	for _, v := range diff.sortedEntries() {
		count[v.Status]++
		if v.Status != Identical {
			fmt.Fprintf(out, "%s %v\n", statusIcon[v.Status], v.File)
		}
	}

	fmt.Fprintf(out,
		"Identical: %v | Modified: %v | Added: %v | Removed: %v\n",
		count[Identical], count[Modified], count[Added], count[Removed])
}

func (diff Diff) sortedEntries() []DiffEntry {
	keys := make([]string, 0, len(diff.Entries))
	values := make([]DiffEntry, 0, len(diff.Entries))

	for k := range diff.Entries {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		values = append(values, diff.Entries[k])
	}

	return values
}
