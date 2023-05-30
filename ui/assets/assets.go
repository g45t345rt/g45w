package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
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
