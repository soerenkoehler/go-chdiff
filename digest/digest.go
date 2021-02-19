package digest

import (
	"fmt"
	"os"
	"path"
	"sort"
	"sync"

	"github.com/soerenkoehler/chdiff-go/util"
)

// Digest is a map file path => checksum
type Digest map[string]string

// Create ... TODO
func Create(dataPath, digestPath, mode string) util.ErrorList {
	digest, errs := calculate(dataPath, mode)
	fmt.Printf("Saving %s\n", digestPath)
	for _, k := range digest.sortedKeys() {
		fmt.Printf("%s => %s\n", k, digest[k])
	}
	return errs
}

// Verify ... TODO
func Verify(dataPath, digestPath, mode string) util.ErrorList {
	digest, errs := calculate(dataPath, mode)
	fmt.Printf("Verify %s\n", digestPath)
	for _, k := range digest.sortedKeys() {
		fmt.Printf("%s => %s\n", k, digest[k])
	}
	return errs
}

func calculate(rootPath, mode string) (Digest, util.ErrorList) {
	errors := util.ErrorList{}

	wait := sync.WaitGroup{}
	pathes := make(chan string)

	addPath := func(path string) {
		wait.Add(1)
		go func() {
			pathes <- path
		}()
	}

	processDir := func(dirPath string) {
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			errors = append(errors, err)
		}
		for _, entry := range entries {
			addPath(path.Join(dirPath, entry.Name()))
		}
	}

	processFile := func(filePath string) {
		fmt.Printf("calculate(%s, %s)\n", mode, filePath)
	}

	go func() {
		for {
			if currentPath, ok := <-pathes; ok {
				info, err := os.Stat(currentPath)
				if err != nil {
					errors = append(errors, err)
				}
				if info.IsDir() {
					processDir(currentPath)
				} else {
					processFile(currentPath)
				}
				wait.Done()
			} else {
				return
			}
		}
	}()

	addPath(rootPath)
	wait.Wait()
	close(pathes)

	return Digest{}, errors
}

func (digest Digest) sortedKeys() []string {
	keys := make([]string, 0, len(digest))
	for key := range digest {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
