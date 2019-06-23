package cache

import (
	"encoding/json"
	"path/filepath"
)

type Cache struct {
	fs  storage
	cfg config
}

type storage interface {
	Exists(name string) bool
	ReadAll(name string) ([]byte, error)
	WriteAll(name string, data []byte) error
}

type config interface {
	FileSystemRoot() string
	DataFormat() string
	EnableCache() bool
}

func New(fs storage, cfg config) *Cache {
	return &Cache{fs, cfg}
}

func (c *Cache) Fetch(
	key string,
	target interface{},
	build func() (interface{}, error),
) error {
	var content []byte
	var err error
	path, err := c.path(key)
	if err != nil {
		return err
	}
	hit := c.fs.Exists(path)
	if c.cfg.EnableCache() && hit {
		return c.fromStorage(path, target)
	}
	data, err := build()
	if err != nil {
		return err
	}
	// Edge-case: if cache is disabled and build() did not overwrite the
	// output file, we get the cached data that could have been written by
	// a preceding successful build instead of an error.
	if data == nil && c.fs.Exists(path) {
		return c.fromStorage(path, target)
	}
	content, err = json.Marshal(data)
	if err != nil {
		return err
	}
	err = c.fs.WriteAll(path, content)
	if err != nil {
		return err
	}
	if c.fs.Exists(path) {
		return c.fromStorage(path, target)
	}
	return nil
}

func (c *Cache) fromStorage(path string, target interface{}) error {
	content, err := c.fs.ReadAll(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, target); err != nil {
		return err
	}
	return nil
}

func (c *Cache) path(key string) (string, error) {
	root := c.cfg.FileSystemRoot()
	file := key + c.cfg.DataFormat()
	if filepath.IsAbs(file) {
		return file, nil
	}
	if filepath.Base(file) == file {
		return filepath.Join(root, file), nil
	}
	rel, err := filepath.Rel(root, file)
	if err != nil {
		return "", err
	}
	return filepath.Join(root, rel), nil
}
