package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"net/http"
)

//go:embed images/*
var images embed.FS

func GetImage(path string) (image.Image, error) {
	data, err := images.ReadFile(fmt.Sprintf("images/%s", path))
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewBuffer(data))
	return img, err
}

func FetchImage(url string) (image.Image, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(res.Body)
	return img, err
}
