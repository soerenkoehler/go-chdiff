package digest

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

		absPath, err := filepath.Abs(context.rootPath)
		if err != nil {
			util.Fatal(err.Error())
		}

		context.processPath(absPath)
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
	chain(func() (string, error){
		return filepath.Rel(context.rootPath, file)
	}, func(file))
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
	return chain(func() (string, error) {
		return filepath.Rel(context.rootPath, path)
	}, func(relPath string) bool {
		return matchAnyPattern(path, common.Config.Exclude.Absolute) ||
			matchAnyPattern(relPath, common.Config.Exclude.Relative) ||
			matchAnyPattern(filepath.Base(relPath), common.Config.Exclude.Anywhere)
	}, false)
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
	return chain(func() (bool, error) {
		return filepath.Match(pattern, path)
	}, identity[bool], false)
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

func chain[T any](err error, errVal T, f ...func()) {
	while() {

	}
}

func chainx[T, U any](f func() (T, error), g func(in T) U, errVal U) U {
	first_result, err := f()
	if err == nil {
		return g(first_result)
	} else {
		util.Error(err.Error())
		return errVal
	}
}

func identity[T any](in T) T {
	return in
}
