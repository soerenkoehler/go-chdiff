package util

import (
	"log"
	"os"
	"path/filepath"
)

// PathInfo distilles information from FileInfo and EvalSymlinks
type PathInfo struct {
	IsDir     bool
	IsSymlink bool
	Target    string
}

// Stat checks if a path is a directory, a symlink or otherwise a regular file.
func Stat(path string) PathInfo {
	info, err := os.Lstat(path)
	if err != nil {
		log.Printf("[E]: %s\n", err)
		return PathInfo{
			IsDir:     false,
			IsSymlink: false,
			Target:    path}
	}
	target, _ := filepath.EvalSymlinks(path)
	return PathInfo{
		IsDir:     info.IsDir(),
		IsSymlink: 0 != (info.Mode() & os.ModeSymlink),
		Target:    target}
}
