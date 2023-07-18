package animation

import (
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
)

type TransformFunc func(gtx layout.Context, value float32) op.TransformOp

func TransformScaleCenter(gtx layout.Context, value float32) op.TransformOp {
	pt := gtx.Constraints.Min.Div(2)
	origin := f32.Pt(float32(pt.X), float32(pt.Y))
	trans := f32.Affine2D{}.Scale(origin, f32.Point{X: value, Y: value})
	return op.Affine(trans)
}

func TransformY(gtx layout.Context, value float32) op.TransformOp {
	pt := f32.Pt(0, float32(gtx.Constraints.Max.Y)*value)
	trans := f32.Affine2D{}.Offset(pt)
	return op.Affine(trans)
}

func TransformX(gtx layout.Context, value float32) op.TransformOp {
	pt := f32.Pt(float32(gtx.Constraints.Max.X)*value, 0)
	trans := f32.Affine2D{}.Offset(pt)
	return op.Affine(trans)
}

func TransformRotate(gtx layout.Context, value float32) op.TransformOp {
	pt := gtx.Constraints.Min.Div(2)
	origin := f32.Pt(float32(pt.X), float32(pt.Y))
	trans := f32.Affine2D{}.Rotate(origin, value*2*math.Pi)
	return op.Affine(trans)
}
