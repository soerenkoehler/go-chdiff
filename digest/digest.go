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
}

// Service is the mockable API for the digest service.
type Service interface {
	Create(dataPath, digestPath, algorithm string) error
	Verify(dataPath, digestPath, algorithm string) error
}

// DefaultService ist the production implementation of the digest service.
type DefaultService struct{}

func (DefaultService) Create(dataPath, digestPath, algorithm string) error {
	digest, err := calculateDigest(dataPath, algorithm)
	fmt.Printf("Saving %s\n", digestPath)
	for _, k := range digest.sortedKeys() {
		fmt.Printf("%s => %v\n", k, digest[k])
	}
	return err
}

func (DefaultService) Verify(dataPath, digestPath, algorithm string) error {
	digest, err := calculateDigest(dataPath, algorithm)
	fmt.Printf("Verify %s\n", digestPath)
	for _, k := range digest.sortedKeys() {
		fmt.Printf("%s => %v\n", k, digest[k])
	}
	return err
}

func calculateDigest(rootpath, algorithm string) (Digest, error) {
	context := DigestContext{
		rootpath:  rootpath,
		algorithm: algorithm,
		waitgroup: &sync.WaitGroup{},
	}
	context.processPath(context.rootpath)
	context.waitgroup.Wait()
	return Digest{}, nil
}

func (context DigestContext) processPath(path string) {
	context.waitgroup.Add(1)
	go func() {
		switch info := util.Stat(path); {
		case info.IsSymlink:
			log.Printf("skipping symlink: %s => %s", path, info.Target)
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
		log.Println(err)
	}
	for _, entry := range entries {
		context.processPath(path.Join(dirPath, entry.Name()))
	}
}

func (context DigestContext) processFile(path string) {
	// fmt.Printf("calculate(%s, %s)\n", context.algorithm, path)
}

func (digest Digest) sortedKeys() []string {
	keys := make([]string, 0, len(digest))
	for key := range digest {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
