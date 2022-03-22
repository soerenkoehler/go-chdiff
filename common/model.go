package common

import "time"

type Location struct {
	Path string
	Time time.Time
}

const LocationTimeFormat = "2006-01-02 15:04:05"

type Set[T comparable] map[T]struct{}

func (set Set[T]) Put(entry T) {
	set[entry] = struct{}{}
}
