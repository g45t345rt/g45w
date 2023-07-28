package components

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type TabBars struct {
	Key     string
	Items   []*TabBarsItem
	changed bool
}

func NewTabBars(defaultKey string, items []*TabBarsItem) *TabBars {
	return &TabBars{
		Key:   defaultKey,
		Items: items,
	}
}

func (t *TabBars) Changed() (bool, string) {
	if t.changed {
		t.changed = false
		return true, t.Key
	}
	return false, t.Key
}

func (t *TabBars) Layout(gtx layout.Context, th *material.Theme, text map[string]string) layout.Dimensions {
	var childrens []layout.FlexChild

	for i, item := range t.Items {
		if item.clickable.Clicked() {
			t.Key = item.Key
			t.changed = true
			op.InvalidateOp{}.Add(gtx.Ops)
		}

		item.selected = t.Key == item.Key
		idx := i
		childrens = append(childrens, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			item := t.Items[idx]
			return item.Layout(gtx, th, text[item.Key])
		}))
	}

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, childrens...)
}

type TabBarItemStyle struct {
	TextSize unit.Sp
}

type TabBarsItem struct {
	Key       string
	clickable *widget.Clickable
	selected  bool
	Style     TabBarItemStyle
}

func NewTabBarItem(key string, style TabBarItemStyle) *TabBarsItem {
	return &TabBarsItem{
		Key:       key,
		selected:  false,
		clickable: new(widget.Clickable),
		Style:     style,
	}
}

func (t *TabBarsItem) Layout(gtx layout.Context, th *material.Theme, text string) layout.Dimensions {
	if t.clickable.Hovered() && !t.selected {
		pointer.CursorPointer.Add(gtx.Ops)
	}

	return t.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top: unit.Dp(10), Bottom: unit.Dp(10),
					Left: unit.Dp(0), Right: unit.Dp(0),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(18), text)
					lbl.Color = color.NRGBA{A: 150}
					lbl.TextSize = t.Style.TextSize

					if t.selected {
						lbl.Color = color.NRGBA{A: 255}
						lbl.Font.Weight = font.Bold
					}

					return lbl.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if t.selected {
					rect := image.Rectangle{Max: image.Pt(gtx.Constraints.Max.X, gtx.Dp(5))}
					paint.FillShape(gtx.Ops,
						color.NRGBA{A: 255},
						clip.UniformRRect(rect, 0).Op(gtx.Ops),
					)
					return layout.Dimensions{Size: rect.Max}
				} else {
					rect := image.Rectangle{Max: image.Pt(gtx.Constraints.Max.X, gtx.Dp(2))}
					paint.FillShape(gtx.Ops,
						color.NRGBA{A: 150},
						clip.UniformRRect(rect, 0).Op(gtx.Ops),
					)
					return layout.Dimensions{Size: rect.Max}
				}
			}),
		)
	})
}
