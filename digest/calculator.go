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

type Calculator func(string, string) Digest

type digestContext struct {
	rootPath  string
	algorithm string
	waitGroup *sync.WaitGroup
	digest    chan digestEntry
}

func Calculate(rootPath, algorithm string) Digest {
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

	result := newDigest(rootPath, algorithm, time.Now())
	for entry := range context.digest {
		result.addEntry(entry)
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
		Hash: hex.EncodeToString(hash.Sum(nil)),
	}
}

func getNewHash(algorithm string) hash.Hash {
	switch algorithm {
	case "SHA256":
		return sha256.New()
	case "SHA512":
		return sha512.New()
	}
	panic(fmt.Errorf("invalid hash algorithm %v", algorithm))
}
