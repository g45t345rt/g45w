package caching

import (
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"

	"github.com/g45t345rt/g45w/settings"
)

func cachePath(relPath string, name string) (cachePath string) {
	cachePath = filepath.Join(settings.CacheDir, relPath, name+".cache")
	return
}

func Store(relPath string, name string, value interface{}) error {
	fullPath := cachePath(relPath, name)

	buffer := &bytes.Buffer{}
	enc := gob.NewEncoder(buffer)
	err := enc.Encode(value)
	if err != nil {
		return err
	}

	dirPath := filepath.Dir(fullPath)

	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	data := buffer.Bytes()
	err = os.WriteFile(fullPath, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func Get(relPath string, name string, value interface{}) (exists bool, err error) {
	fullPath := cachePath(relPath, name)

	_, err = os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		}
		return
	}
	exists = true

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return
	}

	buffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(value)
	if err != nil {
		return
	}

	return
}

func Clear(relCachePath string, name string) error {
	fullPath := cachePath(relCachePath, name)
	dirPath := filepath.Dir(fullPath)
	return os.RemoveAll(dirPath)
}

func ClearAll() error {
	return os.RemoveAll(settings.CacheDir)
}
