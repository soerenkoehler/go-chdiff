package digest

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
	"io/fs"
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

		var absPath string

		chain := &util.ChainContext{}
		chain.Chain(func() {
			absPath, chain.Err = filepath.Abs(context.rootPath)
		}).Chain(func() {
			context.processPath(absPath)
			context.waitGroup.Wait()
		}).ChainFatal("calculate")
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
	var entries []fs.DirEntry

	chain := &util.ChainContext{}
	chain.Chain(func() {
		entries, chain.Err = os.ReadDir(dir)
	}).Chain(func() {
		for _, entry := range entries {
			context.processPath(filepath.Join(dir, entry.Name()))
		}
	}).ChainError("process dir")
}

func (context digestContext) processFile(file string) {
	var relativePath string
	var input *os.File

	chain := &util.ChainContext{}
	chain.Chain(func() {
		relativePath, chain.Err = filepath.Rel(context.rootPath, file)
	}).Chain(func() {
		input, chain.Err = os.Open(file)
	}).Chain(func() {
		defer input.Close()

		hash := getNewHash(context.algorithm)
		io.Copy(hash, input)

		context.digest <- digestEntry{
			file: relativePath,
			hash: hex.EncodeToString(hash.Sum(nil)),
		}
	}).ChainError("process file")
}

func (context digestContext) pathExcluded(path string) bool {
	var relativePath string
	var result bool

	chain := &util.ChainContext{}
	chain.Chain(func() {
		relativePath, chain.Err = filepath.Rel(context.rootPath, path)
	}).Chain(func() {
		result = matchAnyPattern(path, common.Config.Exclude.Absolute) ||
			matchAnyPattern(relativePath, common.Config.Exclude.Relative) ||
			matchAnyPattern(filepath.Base(relativePath), common.Config.Exclude.Anywhere)
	}).ChainError("filter path")

	return result
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
	var chain util.ChainContext
	var result bool

	chain.Chain(func() {
		result, chain.Err = filepath.Match(pattern, path)
	}).ChainError("match path")

	return result
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
