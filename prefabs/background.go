package prefabs

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func PaintGrayLinearGradient(gtx layout.Context) clip.Stack {
	return PaintLinearGradient(gtx,
		color.NRGBA{A: 5},
		color.NRGBA{A: 50},
	)
}

func PaintLinearGradient(gtx layout.Context, colorStart color.NRGBA, colorEnd color.NRGBA) clip.Stack {
	dr := image.Rectangle{Max: gtx.Constraints.Min}
	paint.LinearGradientOp{
		Stop1:  f32.Pt(0, float32(dr.Min.Y)),
		Stop2:  f32.Pt(0, float32(dr.Max.Y)),
		Color1: colorStart,
		Color2: colorEnd,
	}.Add(gtx.Ops)
	stack := clip.Rect(dr).Push(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return stack
}
