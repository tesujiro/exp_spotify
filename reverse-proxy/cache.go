package main

import (
	"encoding/gob"
	"os"
)

type httpCache struct {
	filepath string
	items    map[string][]byte
}

func newHttpCache() httpCache {
	c := httpCache{
		filepath: cacheFilename,
		items:    make(map[string][]byte),
	}
	return c
}

func (c *httpCache) set(key string, value []byte) {
	c.items[key] = value
}

func (c *httpCache) del(key string) {
	delete(c.items, key)
}

func (c *httpCache) get(key string) ([]byte, bool) {
	val, ok := c.items[key]
	return val, ok
}

func (c *httpCache) load() error {
	fd, err := os.Open(c.filepath)
	if err != nil {
		return err
	}
	defer fd.Close()

	dec := gob.NewDecoder(fd)
	err = dec.Decode(&c.items)
	if err != nil {
		return err
	}
	return nil
}

func (c *httpCache) save() error {
	fd, err := os.Create(c.filepath)
	if err != nil {
		return err
	}
	defer fd.Close()

	enc := gob.NewEncoder(fd)
	err = enc.Encode(c.items)
	if err != nil {
		return err
	}
	return nil
}
