// SPDX-License-Identifier: Unlicense OR MIT

package components

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Transform func(dims layout.Dimensions, trans f32.Affine2D) f32.Affine2D

// Image is a widget that displays an image.
type Image struct {
	// Src is the image to display.
	Src paint.ImageOp
	// Fit specifies how to scale the image to the constraints.
	// By default it does not do any scaling.
	Fit Fit
	// Position specifies where to position the image within
	// the constraints.
	Position layout.Direction
	// Scale is the factor used for converting image pixels to dp.
	// If Scale is zero it defaults to 1.
	//
	// To map one image pixel to one output pixel, set Scale to 1.0 / gtx.Metric.PxPerDp.
	Scale float32

	Rounded Rounded
}

func (im Image) Layout(gtx layout.Context, transform Transform) layout.Dimensions {
	scale := im.Scale
	if scale == 0 {
		scale = 1
	}

	size := im.Src.Size()
	wf, hf := float32(size.X), float32(size.Y)
	w, h := gtx.Dp(unit.Dp(wf*scale)), gtx.Dp(unit.Dp(hf*scale))

	//offsetPt := image.Pt(gtx.Dp(100), gtx.Dp(100))

	//constraints := layout.Constraints{
	//	Max: gtx.Constraints.Max.Add(offsetPt.Mul(2)),
	//}

	//defer op.Affine(f32.Affine2D{}.Offset(layout.FPt(offsetPt).Mul(-1))).Push(gtx.Ops).Pop()
	dims, trans := im.Fit.scale(gtx.Constraints, im.Position, layout.Dimensions{Size: image.Pt(w, h)})

	defer clip.RRect{
		Rect: image.Rectangle{Max: dims.Size},
		NW:   gtx.Dp(im.Rounded.NW), NE: gtx.Dp(im.Rounded.NE),
		SE: gtx.Dp(im.Rounded.SE), SW: gtx.Dp(im.Rounded.SW),
	}.Push(gtx.Ops).Pop()

	if transform != nil {
		trans = transform(dims, trans)
	}

	pixelScale := scale * gtx.Metric.PxPerDp
	trans = trans.Mul(f32.Affine2D{}.Scale(f32.Point{}, f32.Pt(pixelScale, pixelScale)))
	defer op.Affine(trans).Push(gtx.Ops).Pop()

	im.Src.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return dims
}
