package cache

import (
	"bytes"
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
}

func New(fs storage, cfg config) *Cache {
	return &Cache{fs, cfg}
}

func (c *Cache) Fetch(
	key string,
	target interface{},
	build func() ([]byte, error),
) error {
	var content []byte
	var err error
	path, err := c.path(key)
	if err != nil {
		return err
	}
	hit := c.fs.Exists(path)
	if hit {
		content, err = c.fs.ReadAll(path)
	} else {
		content, err = build()
	}
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(content)
	err = json.NewDecoder(b).Decode(target)
	if err != nil {
		return err
	}
	if !hit {
		err = c.fs.WriteAll(path, content)
	}
	if err != nil {
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
	rel, err := filepath.Rel(root, file)
	if err != nil {
		return "", err
	}
	return filepath.Join(root, rel), nil
}

// TODO: improve the way we can enforce the storage format.
func (Cache) Serialize(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
