package cache

import (
	"bytes"
	"encoding/json"
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
	path := key + c.cfg.DataFormat()
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

// TODO: improve the way we can enforce the storage format.
func (Cache) Serialize(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
