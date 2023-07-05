package prefabs

import (
	"fmt"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/components"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type RingSizeSelector struct {
	buttonSelect *components.Button
	selectModal  *SelectModal

	changed bool
	Value   string
}

func NewRingSizeSelector(defaultSize string) *RingSizeSelector {
	tuneIcon, _ := widget.NewIcon(icons.ActionTrackChanges)
	buttonSelect := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
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

	sizes := []string{"2", "4", "8", "16", "32", "64", "128"}
	items := []*SelectListItem{}

	for _, size := range sizes {
		items = append(items, NewSelectListItem(size, func(gtx layout.Context, index int, th *material.Theme) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(20), sizes[index])
			return lbl.Layout(gtx)
		}))
	}

	selectModal := NewSelectModal()
	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			selectModal.Layout(gtx, th, items)
		},
	})

	r := &RingSizeSelector{
		buttonSelect: buttonSelect,
		selectModal:  selectModal,
		Value:        defaultSize,
	}

	return r
}

func (r *RingSizeSelector) Changed() bool {
	return r.changed
}

func (r *RingSizeSelector) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	r.changed = false

	if r.buttonSelect.Clickable.Clicked() {
		r.selectModal.Modal.SetVisible(gtx, true)
	}

	selected := r.selectModal.Selected()
	if selected {
		r.Value = r.selectModal.SelectedKey
		r.changed = true
		r.selectModal.Modal.SetVisible(gtx, false)
	}

	r.buttonSelect.Text = fmt.Sprintf("Ring size: %s", r.Value)
	return r.buttonSelect.Layout(gtx, th)
}
