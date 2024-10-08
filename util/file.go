package util

import (
	"os"
	"path/filepath"
)

// PathInfo distilles information from FileInfo and Readlink
type PathInfo struct {
	IsDir     bool
	IsSymlink bool
	Target    string
}

// Stat checks if a path is a directory, a symlink or otherwise a regular file.
func Stat(path string) PathInfo {
	defaultResult := PathInfo{
		IsDir:     false,
		IsSymlink: false,
		Target:    path}

	info, err := os.Lstat(path)
	if err != nil {
		Error(err.Error())
		return defaultResult
	}

	target, err := filepath.EvalSymlinks(path)
	if err != nil {
		Error(err.Error())
		return defaultResult
	}

	return PathInfo{
		IsDir:     info.IsDir(),
		IsSymlink: (info.Mode() & os.ModeSymlink) != 0,
		Target:    target}
}
