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

type TabBarsColors struct {
	ActiveColor   color.NRGBA
	InactiveColor color.NRGBA
}

type TabBars struct {
	Key    string
	Items  []*TabBarsItem
	Colors TabBarsColors

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

func (t *TabBars) Layout(gtx layout.Context, th *material.Theme, textSize unit.Sp, text map[string]string) layout.Dimensions {
	var childrens []layout.FlexChild

	for i, item := range t.Items {
		if item.clickable.Clicked(gtx) {
			t.Key = item.Key
			t.changed = true
			op.InvalidateOp{}.Add(gtx.Ops)
		}

		item.selected = t.Key == item.Key
		idx := i
		childrens = append(childrens, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			item := t.Items[idx]
			value := text[item.Key]
			return item.Layout(gtx, th, t, textSize, value)
		}))
	}

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, childrens...)
}

type TabBarsItem struct {
	Key       string
	clickable *widget.Clickable
	selected  bool
}

func NewTabBarItem(key string) *TabBarsItem {
	return &TabBarsItem{
		Key:       key,
		selected:  false,
		clickable: new(widget.Clickable),
	}
}

func (t *TabBarsItem) Layout(gtx layout.Context, th *material.Theme, tabBars *TabBars, textSize unit.Sp, text string) layout.Dimensions {
	return t.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if t.clickable.Hovered() && !t.selected {
			pointer.CursorPointer.Add(gtx.Ops)
		}

		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top: unit.Dp(10), Bottom: unit.Dp(10),
					Left: unit.Dp(0), Right: unit.Dp(0),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, textSize, text)
					lbl.Color = tabBars.Colors.InactiveColor

					if t.selected {
						lbl.Color = tabBars.Colors.ActiveColor
						lbl.Font.Weight = font.Bold
					}

					return lbl.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if t.selected {
					rect := image.Rectangle{Max: image.Pt(gtx.Constraints.Max.X, gtx.Dp(5))}
					paint.FillShape(gtx.Ops,
						tabBars.Colors.ActiveColor,
						clip.UniformRRect(rect, 0).Op(gtx.Ops),
					)
					return layout.Dimensions{Size: rect.Max}
				} else {
					rect := image.Rectangle{Max: image.Pt(gtx.Constraints.Max.X, gtx.Dp(2))}
					paint.FillShape(gtx.Ops,
						tabBars.Colors.InactiveColor,
						clip.UniformRRect(rect, 0).Op(gtx.Ops),
					)
					return layout.Dimensions{Size: rect.Max}
				}
			}),
		)
	})
}
