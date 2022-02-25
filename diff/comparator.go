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

var statusIcon map[diffStatus]string = map[diffStatus]string{
	Identical: " ",
	Modified:  "*",
	Added:     "+",
	Removed:   "-"}

func Compare(old, new digest.Digest) Diff {
	return Diff{
		locationA: old.Location,
		locationB: new.Location,
		entries:   diffEntries{}}
}

func Print(out io.Writer, diff Diff) {
	fmt.Fprintf(out,
		"Old: (%s) %v\nNew: (%s) %v\n",
		diff.locationA.Time.Format(common.LocationTimeFormat),
		diff.locationA.Path,
		diff.locationB.Time.Format(common.LocationTimeFormat),
		diff.locationB.Path)

	count := make(map[diffStatus]int32, 4)

	for _, v := range diff.sortedEntries() {
		count[v.status]++
		if v.status != Identical {
			fmt.Fprintf(out, "%s %v\n", statusIcon[v.status], v.file)
		}
	}

	fmt.Fprintf(out,
		"Identical: %v | Modified: %v | Added: %v | Removed: %v\n",
		count[Identical], count[Modified], count[Added], count[Removed])
}

func (diff Diff) sortedEntries() []diffEntry {
	keys := make([]string, 0, len(diff.entries))
	values := make([]diffEntry, 0, len(diff.entries))

	for k := range diff.entries {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		values = append(values, diff.entries[k])
	}

	return values
}
