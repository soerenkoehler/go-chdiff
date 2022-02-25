package util

import (
	"log"
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
	info, err := os.Lstat(path)
	if err == nil {
		target, err := filepath.EvalSymlinks(path)
		if err == nil {
			return PathInfo{
				IsDir:     info.IsDir(),
				IsSymlink: (info.Mode() & os.ModeSymlink) != 0,
				Target:    target}
		}
	}
	log.Printf("[E]: %s\n", err)
	return PathInfo{
		IsDir:     false,
		IsSymlink: false,
		Target:    path}
}
