package prefabs

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/g45t345rt/g45w/theme"
)

func Divider(gtx layout.Context, height unit.Dp) layout.Dimensions {
	gtx.Constraints.Max.Y = gtx.Dp(height)
	paint.FillShape(gtx.Ops, theme.Current.DividerColor,
		clip.Rect{
			Max: gtx.Constraints.Max,
		}.Op())
	return layout.Dimensions{Size: gtx.Constraints.Max}
}
