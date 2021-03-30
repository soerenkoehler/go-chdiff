package util

import (
	"io"
	"os"
)

// StdIOService provides mockable console IO.
type StdIOService interface {
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
}

// DefaultStdIOService provides the default IO from the os package.
type DefaultStdIOService struct{}

// Stdin provides the default stdin.
func (DefaultStdIOService) Stdin() io.Reader {
	return os.Stdin
}

// Stdout provides the default stdout.
func (DefaultStdIOService) Stdout() io.Writer {
	return os.Stdout
}

// Stderr provides the default stderr.
func (DefaultStdIOService) Stderr() io.Writer {
	return os.Stderr
}
