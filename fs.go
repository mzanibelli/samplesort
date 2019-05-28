package samplesort

import (
	"io"
	"os"
	"path/filepath"
)

// See: https://stackoverflow.com/a/16742530

var fs filesystem = osFS{}

type filesystem interface {
	Open(name string) (file, error)
	Stat(name string) (os.FileInfo, error)
	Exists(name string) bool
	Walk(name string, f filepath.WalkFunc) error
}

type file interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
}

type osFS struct{}

func (osFS) Open(name string) (file, error) { return os.Open(name) }

func (osFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }

func (osFS) Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func (osFS) Walk(name string, f filepath.WalkFunc) error {
	return filepath.Walk(name, f)
}
