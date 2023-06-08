package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type ProgressBarStyle struct {
}

type ProgressBar struct {
	Value float32

	Color   color.NRGBA
	BgColor color.NRGBA
	Rounded unit.Dp
	Height  unit.Dp
}

func (p ProgressBar) Layout(gtx layout.Context) layout.Dimensions {
	rounded := gtx.Dp(p.Rounded)

	size := gtx.Constraints.Min.Add(image.Pt(0, gtx.Dp(p.Height)))

	defer clip.RRect{
		Rect: image.Rectangle{
			Max: size,
		},
		SE: rounded, SW: rounded,
		NW: rounded, NE: rounded,
	}.Push(gtx.Ops).Pop()

	paint.ColorOp{Color: p.BgColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	//paint.FillShape(gtx.Ops, p.BgColor, .Op(gtx.Ops))

	valueWidth := unit.Dp(p.Value * float32(gtx.Constraints.Min.X))
	defer clip.Rect{Max: image.Pt(gtx.Dp(valueWidth), gtx.Dp(p.Height))}.Push(gtx.Ops).Pop()

	paint.ColorOp{Color: p.Color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: size}
}
