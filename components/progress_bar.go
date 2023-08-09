package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type ProgressBarColors struct {
	IndicatorColor  color.NRGBA
	BackgroundColor color.NRGBA
}

type ProgressBar struct {
	Value   float32
	Rounded unit.Dp
	Height  unit.Dp
	Colors  ProgressBarColors
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

	backgroundColor := p.Colors.BackgroundColor
	paint.ColorOp{Color: backgroundColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	//paint.FillShape(gtx.Ops, p.BgColor, .Op(gtx.Ops))

	valueWidth := unit.Dp(p.Value * float32(gtx.Constraints.Min.X))
	defer clip.Rect{Max: image.Pt(gtx.Dp(valueWidth), gtx.Dp(p.Height))}.Push(gtx.Ops).Pop()

	indicatorColor := p.Colors.IndicatorColor
	paint.ColorOp{Color: indicatorColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: size}
}
