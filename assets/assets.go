package assets

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"image"
	"sync"
)

//go:embed images/*
var images embed.FS
var imageMemCache map[string]image.Image
var imageMemCacheMutex sync.Mutex

func GetImage(path string) (image.Image, error) {
	if imageMemCache == nil {
		imageMemCacheMutex.Lock()
		imageMemCache = make(map[string]image.Image)
		imageMemCacheMutex.Unlock()
	}

	img, ok := imageMemCache[path]
	if ok {
		return img, nil
	}

	data, err := images.ReadFile(fmt.Sprintf("images/%s", path))
	if err != nil {
		return nil, err
	}

	newImg, _, err := image.Decode(bytes.NewBuffer(data))
	imageMemCacheMutex.Lock()
	imageMemCache[path] = newImg
	imageMemCacheMutex.Unlock()

	return newImg, err
}

//go:embed lang/*
var lang embed.FS

func GetLang(path string) (map[string]string, error) {
	data, err := lang.ReadFile(fmt.Sprintf("lang/%s", path))
	if err != nil {
		return nil, err
	}

	var values map[string]string
	err = json.Unmarshal(data, &values)
	if err != nil {
		return nil, err
	}

	return values, err
}

//go:embed fonts/*
var font embed.FS

func GetFont(path string) ([]byte, error) {
	data, err := font.ReadFile(fmt.Sprintf("fonts/%s", path))
	if err != nil {
		return nil, err
	}

	return data, err
}
