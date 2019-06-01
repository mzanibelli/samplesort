package cache

import (
	"bytes"
	"encoding/json"
)

type Cache struct {
	fs     storage
	format string
}

type storage interface {
	Exists(name string) bool
	ReadAll(name string) ([]byte, error)
	WriteAll(name string, data []byte) error
}

func New(fs storage, format string) *Cache {
	return &Cache{fs, format}
}

func (c *Cache) Fetch(
	key string,
	target interface{},
	build func() ([]byte, error),
) error {
	var content []byte
	var err error
	path := c.path(key)
	warm := c.fs.Exists(path)
	if warm {
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
	if !warm {
		err = c.fs.WriteAll(path, content)
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) path(key string) string {
	return key + c.format
}
