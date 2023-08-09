package prefabs

import (
	"image"
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/theme"
)

type SelectModal struct {
	Modal *components.Modal

	list        *widget.List
	selected    bool
	selectedKey string
}

func NewSelectModal() *SelectModal {
	list := new(widget.List)
	list.Axis = layout.Vertical

	modal := components.NewModal(components.ModalStyle{
		CloseOnOutsideClick: true,
		CloseOnInsideClick:  false,
		Direction:           layout.S,
		Rounded:             components.UniformRounded(unit.Dp(10)),
		Inset:               layout.UniformInset(25),
		Animation:           components.NewModalAnimationUp(),
	})

	return &SelectModal{
		Modal: modal,
		list:  list,
	}
}

func (l *SelectModal) Selected() (bool, string) {
	return l.selected, l.selectedKey
}

func (l *SelectModal) Layout(gtx layout.Context, th *material.Theme, items []*SelectListItem) layout.Dimensions {
	l.selected = false
	l.Modal.Style.Colors = theme.Current.ModalColors
	return l.Modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(10), Bottom: unit.Dp(10),
			Left: unit.Dp(10), Right: unit.Dp(10),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.Y = gtx.Dp(200)
			listStyle := material.List(th, l.list)
			listStyle.AnchorStrategy = material.Overlay

			return listStyle.Layout(gtx, len(items), func(gtx layout.Context, index int) layout.Dimensions {
				if items[index].clickable.Clicked() {
					l.selectedKey = items[index].Key
					l.selected = true
					op.InvalidateOp{}.Add(gtx.Ops)
				}

				return items[index].Layout(gtx, th, index)
			})
		})
	})
}

type SelectListItem struct {
	Key       string
	element   SelectListItemElement
	clickable *widget.Clickable
}

type SelectListItemElement = func(gtx layout.Context, index int, th *material.Theme) layout.Dimensions

func NewSelectListItem(key string, element SelectListItemElement) *SelectListItem {
	return &SelectListItem{
		Key:       key,
		element:   element,
		clickable: new(widget.Clickable),
	}
}

func (c *SelectListItem) Layout(gtx layout.Context, th *material.Theme, index int) layout.Dimensions {
	dims := c.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return c.element(gtx, index, th)
		})
	})

	if c.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)

		paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 100},
			clip.UniformRRect(
				image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)},
				gtx.Dp(15),
			).Op(gtx.Ops),
		)
	}

	return dims
}
