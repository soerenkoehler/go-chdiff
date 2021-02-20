package digest

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"sync"

	"github.com/soerenkoehler/chdiff-go/util"
)

// Digest is a map file path => checksum
type Digest map[string]string

// Create ... TODO
func Create(dataPath, digestPath, mode string) error {
	digest, err := calculate(dataPath, mode)
	fmt.Printf("Saving %s\n", digestPath)
	for _, k := range digest.sortedKeys() {
		fmt.Printf("%s => %s\n", k, digest[k])
	}
	return err
}

// Verify ... TODO
func Verify(dataPath, digestPath, mode string) error {
	digest, err := calculate(dataPath, mode)
	fmt.Printf("Verify %s\n", digestPath)
	for _, k := range digest.sortedKeys() {
		fmt.Printf("%s => %s\n", k, digest[k])
	}
	return err
}

func calculate(rootPath, mode string) (Digest, error) {
	wait := sync.WaitGroup{}

	var processPath, processDir, processFile func(string)

	processPath = func(path string) {
		wait.Add(1)
		go func() {
			switch info := util.Stat(path); {
			case info.IsSymlink:
				log.Printf("skipping symlink: %s => %s", path, info.Target)
			case info.IsDir:
				processDir(path)
			default:
				processFile(path)
			}
			wait.Done()
		}()
	}

	processDir = func(dirPath string) {
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			log.Println(err)
		}
		for _, entry := range entries {
			processPath(path.Join(dirPath, entry.Name()))
		}
	}

	processFile = func(path string) {
		// fmt.Printf("calculate(%s, %s)\n", mode, path)
	}

	processPath(rootPath)
	wait.Wait()

	return Digest{}, nil
}

func (digest Digest) sortedKeys() []string {
	keys := make([]string, 0, len(digest))
	for key := range digest {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
