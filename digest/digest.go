package digest

import "fmt"

// Create ... TODO
func Create(dataPath, digestPath, mode string) []error {
	errs := calculate(dataPath, mode)
	fmt.Printf("Saving %s\n", digestPath)
	return errs
}

// Verify ... TODO
func Verify(dataPath, digestPath, mode string) []error {
	errs := calculate(dataPath, mode)
	fmt.Printf("Verify %s\n", digestPath)
	return errs
}

func calculate(dataPath, mode string) []error {
	fmt.Printf("Calculate %s Digest for %s\n", mode, dataPath)
	return []error{}
}
