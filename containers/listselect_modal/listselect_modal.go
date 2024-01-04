package listselect_modal

import (
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
)

type ThemeWidget = func(gtx layout.Context, th *material.Theme) layout.Dimensions

type ListSelectModal struct {
	Modal *components.Modal

	list    *widget.List
	keyChan chan string
	items   []*SelectListItem
}

var Instance *ListSelectModal

func LoadInstance() {
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

	Instance = &ListSelectModal{
		Modal: modal,
		list:  list,
		items: make([]*SelectListItem, 0),
	}

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 2,
		Layout: func(gtx layout.Context, th *material.Theme) {
			Instance.Layout(gtx, th)
		},
	})
}

func (l *ListSelectModal) Open(items []*SelectListItem) chan string {
	l.items = items
	l.Modal.SetVisible(true)
	l.keyChan = make(chan string)
	return l.keyChan
}

func (l *ListSelectModal) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	l.Modal.Style.Colors = theme.Current.ModalColors
	return l.Modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(10), Bottom: unit.Dp(10),
			Left: unit.Dp(10), Right: unit.Dp(10),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.Y = gtx.Dp(200)
			listStyle := material.List(th, l.list)
			listStyle.AnchorStrategy = material.Overlay

			return listStyle.Layout(gtx, len(l.items), func(gtx layout.Context, index int) layout.Dimensions {
				if l.items[index].clickable.Clicked(gtx) {
					l.keyChan <- l.items[index].Key
					l.Modal.SetVisible(false)
					//close(l.keyChan) don't close channel it panics if spamming click
					op.InvalidateOp{}.Add(gtx.Ops)
				}

				return l.items[index].Layout(gtx, th)
			})
		})
	})
}

type SelectListItem struct {
	Key       string
	element   ThemeWidget
	clickable *widget.Clickable
}

func NewSelectListItem(key string, element ThemeWidget) *SelectListItem {
	return &SelectListItem{
		Key:       key,
		element:   element,
		clickable: new(widget.Clickable),
	}
}

func (item *SelectListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	r := op.Record(gtx.Ops)
	dims := item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return item.element(gtx, th)
		})
	})
	c := r.Stop()

	if item.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)

		paint.FillShape(gtx.Ops, theme.Current.ListItemHoverBgColor,
			clip.UniformRRect(
				image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)},
				gtx.Dp(15),
			).Op(gtx.Ops),
		)
	}

	c.Add(gtx.Ops)
	return dims
}

type ItemText struct {
	icon *widget.Icon
	text string
}

func NewItemText(icon *widget.Icon, text string) *ItemText {
	return &ItemText{
		icon: icon,
		text: text,
	}
}

func (item *ItemText) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	var childs []layout.FlexChild

	if item.icon != nil {
		childs = append(childs,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return item.icon.Layout(gtx, th.Fg)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
		)
	}

	childs = append(childs, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(20), item.text)
		return lbl.Layout(gtx)
	}))

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx, childs...)
}
