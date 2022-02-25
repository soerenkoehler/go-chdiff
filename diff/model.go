package diff

import "github.com/soerenkoehler/go-chdiff/common"

type diffStatus int32

const (
	Identical diffStatus = iota
	Modified
	Added
	Removed
)

type diffEntry struct {
	file   string
	status diffStatus
}

type diffEntries map[string]diffEntry

type Diff struct {
	locationA common.Location
	locationB common.Location
	entries   diffEntries
}
