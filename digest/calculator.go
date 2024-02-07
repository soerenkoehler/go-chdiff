package digest

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/soerenkoehler/go-chdiff/util"
)

type Calculator func(
	rootPath string,
	exclude util.PathFilter,
	algorithm HashType) Digest

type digestEntry struct {
	file string
	hash string
}

type digestContext struct {
	rootPath  string
	exclude   util.PathFilter
	algorithm HashType
	waitGroup *sync.WaitGroup
	digest    chan digestEntry
}

func Calculate(
	rootPath string,
	exclude util.PathFilter,
	algorithm HashType) Digest {

	context := digestContext{
		rootPath:  rootPath,
		exclude:   exclude,
		algorithm: algorithm,
		waitGroup: &sync.WaitGroup{},
		digest:    make(chan digestEntry),
	}

	go func() {
		defer close(context.digest)

		absPath, err := filepath.Abs(context.rootPath)
		if err == nil {
			return
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

		if context.exclude.Matches(path) {
			return
		}

		switch info := util.Stat(path); {
		case info.IsSymlink:
			log.Printf("[W] skipping symlink: %s => %s", path, info.Target)
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
		log.Printf("[E]: %s\n", err)
		return
	}

	for _, entry := range entries {
		context.processPath(filepath.Join(dir, entry.Name()))
	}
}

func (context digestContext) processFile(file string) {
	relativePath, err := filepath.Rel(context.rootPath, file)
	if err != nil {
		log.Printf("[E]: %s\n", err)
		return
	}

	input, err := os.Open(file)
	if err != nil {
		log.Printf("[E]: %s\n", err)
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
