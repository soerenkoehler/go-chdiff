package digest

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/soerenkoehler/go-chdiff/common"
	"github.com/soerenkoehler/go-chdiff/util"
)

type Calculator func(
	rootPath string,
	algorithm HashType) Digest

type digestEntry struct {
	file string
	hash string
}

type digestContext struct {
	rootPath  string
	algorithm HashType
	waitGroup *sync.WaitGroup
	digest    chan digestEntry
}

func Calculate(
	rootPath string,
	algorithm HashType) Digest {

	context := digestContext{
		rootPath:  rootPath,
		algorithm: algorithm,
		waitGroup: &sync.WaitGroup{},
		digest:    make(chan digestEntry),
	}

	go func() {
		defer close(context.digest)

		context.processPath(".")
		context.waitGroup.Wait()
	}()

	result := NewDigest(rootPath, time.Now())
	for entry := range context.digest {
		(*result.Entries)[entry.file] = entry.hash
	}

	return result
}

func (context digestContext) processPath(path string) {
	context.waitGroup.Add(1)
	go func() {
		defer context.waitGroup.Done()

		if context.pathExcluded(path) {
			return
		}

		switch info := util.Stat(path); {
		case info.IsSymlink:
			util.Warn("skipping symlink: %v -> %v", path, info.Target)
		case info.IsDir:
			context.processDir(path)
		default:
			context.processFile(path)
		}
	}()
}

func (context digestContext) processDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		util.Error(err.Error())
		return
	}

	for _, entry := range entries {
		context.processPath(filepath.Join(dir, entry.Name()))
	}
}

func (context digestContext) processFile(file string) {
	relativePath, err := filepath.Rel(context.rootPath, file)
	if err != nil {
		util.Error(err.Error())
		return
	}

	input, err := os.Open(file)
	if err != nil {
		util.Error(err.Error())
		return
	}

	defer input.Close()

	hash := getNewHash(context.algorithm)
	io.Copy(hash, input)

	context.digest <- digestEntry{
		file: relativePath,
		hash: hex.EncodeToString(hash.Sum(nil)),
	}
}

func (context digestContext) pathExcluded(path string) bool {
	absPath, err := filepath.Abs(filepath.Join(context.rootPath, path))
	if err != nil {
		util.Error(err.Error())
	}

	return matchAnyPattern(absPath, common.Config.Exclude.Absolute) ||
		matchAnyPattern(path, common.Config.Exclude.RootRelative) ||
		matchAnyPattern(filepath.Base(path), common.Config.Exclude.Anywhere)
}

func matchAnyPattern(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if matchPattern(pattern, path) {
			return true
		}
	}
	return false
}

func matchPattern(path, pattern string) bool {
	match, err := filepath.Match(pattern, path)
	if err != nil {
		util.Error(err.Error())
		return false
	}
	return match
}

func getNewHash(algorithm HashType) hash.Hash {
	switch algorithm {
	case SHA512:
		return sha512.New()
	case SHA256:
		fallthrough
	default:
		return sha256.New()
	}
}
