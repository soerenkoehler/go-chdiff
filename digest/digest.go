package digest

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"sync"
	"time"

	"github.com/soerenkoehler/chdiff-go/util"
)

type DigestEntry struct {
	path    string
	hash    string
	size    int64
	modTime time.Time
}

type Digest map[string]DigestEntry

type DigestContext struct {
	rootpath  string
	algorithm string
	waitgroup *sync.WaitGroup
	digest    chan DigestEntry
}

// Service is the mockable API for the digest service.
type Service interface {
	Create(dataPath, digestPath, algorithm string)
	Verify(dataPath, digestPath, algorithm string)
}

// DefaultService ist the production implementation of the digest service.
type DefaultService struct{}

func (DefaultService) Create(dataPath, digestPath, algorithm string) {
	digest := calculateDigest(dataPath, algorithm)
	fmt.Printf("Saving %s\n", digestPath)
	for _, k := range digest.sortedKeys() {
		fmt.Print(digest[k].entryToString())
	}
}

func (DefaultService) Verify(dataPath, digestPath, algorithm string) {
	digest := calculateDigest(dataPath, algorithm)
	fmt.Printf("Verify %s\n", digestPath)
	for _, k := range digest.sortedKeys() {
		fmt.Print(digest[k].entryToString())
	}
}

func calculateDigest(rootpath, algorithm string) Digest {
	context := DigestContext{
		rootpath:  rootpath,
		algorithm: algorithm,
		waitgroup: &sync.WaitGroup{},
		digest:    make(chan DigestEntry),
	}

	go func() {
		context.processPath(context.rootpath)
		context.waitgroup.Wait()
		close(context.digest)
	}()

	result := Digest{}
	for entry := range context.digest {
		result[entry.path] = entry
	}

	return result
}

func (context DigestContext) processPath(path string) {
	context.waitgroup.Add(1)
	go func() {
		switch info := util.Stat(path); {
		case info.IsSymlink:
			log.Printf("[W] skipping symlink: %s => %s", path, info.Target)
		case info.IsDir:
			context.processDir(path)
		default:
			context.processFile(path)
		}
		context.waitgroup.Done()
	}()
}

func (context DigestContext) processDir(dirPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Printf("[E]: %s\n", err)
	} else {
		for _, entry := range entries {
			context.processPath(path.Join(dirPath, entry.Name()))
		}
	}
}

func (context DigestContext) processFile(path string) {
	info, err := os.Lstat(path)
	if err != nil {
		log.Printf("[E]: %s\n", err)
	} else {
		context.digest <- DigestEntry{
			path:    path,
			hash:    "tbd",
			size:    info.Size(),
			modTime: info.ModTime(),
		}
	}
}

func (digest Digest) sortedKeys() []string {
	keys := make([]string, 0, len(digest))
	for key := range digest {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (entry DigestEntry) entryToString() string {
	return fmt.Sprintf(
		"# %d %s %s\n%s *%s\n",
		entry.size,
		entry.modTime.Local().Format("20060102-150405"),
		entry.path,
		entry.hash,
		entry.path)
}
