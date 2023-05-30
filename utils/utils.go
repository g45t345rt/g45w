package utils

import (
	"image"
	"image/color"

	"gioui.org/gpu/headless"
	"gioui.org/layout"
)

// From gio material theme
func HexColor(c uint32) color.NRGBA {
	return RGBA(0xff000000 | c)
}

func RGBA(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func CaptureFrame(gtx layout.Context, img *image.RGBA) error {
	rect := image.Rectangle{Max: gtx.Constraints.Max}
	w, err := headless.NewWindow(rect.Max.X, rect.Max.Y)
	if err != nil {
		return err
	}

	w.Frame(gtx.Ops)
	err = w.Screenshot(img)
	if err != nil {
		return err
	}

	return nil
}

func NewImageColor(size image.Point, color color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			img.SetRGBA(x, y, color)
		}
	}

	return img
}
