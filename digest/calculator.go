package digest

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/soerenkoehler/go-chdiff/util"
)

type Calculator func(rootPath string, algorithm HashType) Digest

// TODO locale type like digestEntry

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

func Calculate(rootPath string, algorithm HashType) Digest {
	context := digestContext{
		rootPath:  rootPath,
		algorithm: algorithm,
		waitGroup: &sync.WaitGroup{},
		digest:    make(chan digestEntry),
	}

	go func() {
		context.processPath(context.rootPath)
		context.waitGroup.Wait()
		close(context.digest)
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
		switch info := util.Stat(path); {
		case info.IsSymlink:
			log.Printf("[W] skipping symlink: %s => %s", path, info.Target)
		case info.IsDir:
			context.processDir(path)
		default:
			context.processFile(path)
		}
		context.waitGroup.Done()
	}()
}

func (context digestContext) processDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("[E]: %s\n", err)
		return
	}

	for _, entry := range entries {
		context.processPath(path.Join(dir, entry.Name()))
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
	case SHA256:
		return sha256.New()
	case SHA512:
		return sha512.New()
	}
	panic(fmt.Errorf("invalid hash algorithm %v", algorithm))
}
