package prefabs

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/ui/components"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type RingSizeSelector struct {
	buttonSelect  *components.Button
	listSizeModal *ListSizeModal

	changed bool
	Value   string
}

func NewRingSizeSelector(defaultSize string) *RingSizeSelector {
	tuneIcon, _ := widget.NewIcon(icons.ActionTrackChanges)
	buttonSelect := components.NewButton(components.ButtonStyle{
		Rounded:         unit.Dp(5),
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(16),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Icon:            tuneIcon,
		IconGap:         unit.Dp(10),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonSelect.Label.Alignment = text.Middle
	buttonSelect.Style.Font.Weight = font.Bold

	sizes := []string{"2", "4", "8", "16", "32", "64", "128", "256"}
	items := []*ListSizeItem{}

	for _, size := range sizes {
		items = append(items, NewListSizeItem(size))
	}

	w := app_instance.Current.Window
	th := app_instance.Current.Theme
	router := app_instance.Current.Router
	listSizeModal := NewListSizeModal(w, th)
	router.PushLayout(func(gtx layout.Context, th *material.Theme) {
		listSizeModal.Layout(gtx, th, items)
	})

	r := &RingSizeSelector{
		buttonSelect:  buttonSelect,
		listSizeModal: listSizeModal,
		Value:         defaultSize,
	}

	r.setButtonText(defaultSize)
	return r
}
func (r *RingSizeSelector) setButtonText(value string) {
	r.buttonSelect.Style.Text = fmt.Sprintf("Ring size: %s", value)
}

func (r *RingSizeSelector) Changed() bool {
	if r.changed {
		r.changed = false
		return true
	}

	return false
}

func (r *RingSizeSelector) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if r.buttonSelect.Clickable.Clicked() {
		r.listSizeModal.modal.SetVisible(true)
	}

	changed := r.listSizeModal.Changed()
	if changed {
		r.Value = r.listSizeModal.Value
		r.setButtonText(r.Value)
		r.changed = true
		r.listSizeModal.modal.SetVisible(false)
	}

	return r.buttonSelect.Layout(gtx, th)
}

type ListSizeModal struct {
	listStyle material.ListStyle
	modal     *components.Modal

	changed bool
	Value   string
}

func NewListSizeModal(w *app.Window, th *material.Theme) *ListSizeModal {
	list := new(widget.List)
	list.Axis = layout.Vertical
	listStyle := material.List(th, list)
	listStyle.AnchorStrategy = material.Overlay

	modal := components.NewModal(w, components.ModalStyle{
		CloseOnOutsideClick: true,
		CloseOnInsideClick:  false,
		Direction:           layout.S,
		BgColor:             color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Rounded:             unit.Dp(10),
		Inset:               layout.UniformInset(25),
		Animation:           components.NewModalAnimationUp(),
		Backdrop:            components.NewModalBackground(),
	})

	return &ListSizeModal{
		modal:     modal,
		listStyle: listStyle,
	}
}

func (l *ListSizeModal) Changed() bool {
	if l.changed {
		l.changed = false
		return true
	}

	return false
}

func (l *ListSizeModal) Layout(gtx layout.Context, th *material.Theme, items []*ListSizeItem) layout.Dimensions {
	return l.modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(10), Bottom: unit.Dp(10),
			Left: unit.Dp(10), Right: unit.Dp(10),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return l.listStyle.Layout(gtx, len(items), func(gtx layout.Context, index int) layout.Dimensions {
				if items[index].clickable.Clicked() {
					l.Value = items[index].value
					l.changed = true
				}

				return items[index].Layout(gtx, th)
			})
		})
	})
}

type ListSizeItem struct {
	value     string
	clickable *widget.Clickable
}

func NewListSizeItem(value string) *ListSizeItem {
	return &ListSizeItem{
		value:     value,
		clickable: new(widget.Clickable),
	}
}

func (c *ListSizeItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	dims := c.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			label := material.Label(th, unit.Sp(20), c.value)
			return label.Layout(gtx)
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
