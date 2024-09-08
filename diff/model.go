package diff

import "github.com/soerenkoehler/go-chdiff/common"

type DiffStatus int32

const (
	Identical DiffStatus = iota
	Modified
	Added
	Removed
)

type DiffEntry struct {
	File   string
	Status DiffStatus
}

type Diff struct {
	LocationA common.Location
	LocationB common.Location
	Entries   map[string]DiffEntry
}
