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

	var processPath, processDir, processFile func(string)

	// TODO skip symbolic links
	processPath = func(path string) {
		wait.Add(1)
		go func() {
			info, err := os.Stat(path)
			if err != nil {
				errors = append(errors, err)
			} else {
				if info.IsDir() {
					processDir(path)
				} else {
					processFile(path)
				}
			}
			wait.Done()
		}()
	}

	processDir = func(dirPath string) {
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			errors = append(errors, err)
		}
		for _, entry := range entries {
			processPath(path.Join(dirPath, entry.Name()))
		}
	}

	processFile = func(path string) {
		fmt.Printf("calculate(%s, %s)\n", mode, path)
	}

	processPath(rootPath)
	wait.Wait()

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
