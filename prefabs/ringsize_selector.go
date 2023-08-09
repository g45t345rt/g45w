package prefabs

import (
	"fmt"
	"strconv"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type RingSizeSelector struct {
	buttonSelect *components.Button
	selectModal  *SelectModal

	changed bool
	Value   int
}

func NewRingSizeSelector(defaultSize int) *RingSizeSelector {
	tuneIcon, _ := widget.NewIcon(icons.ActionTrackChanges)
	buttonSelect := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		TextSize:  unit.Sp(16),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Icon:      tuneIcon,
		IconGap:   unit.Dp(10),
		Animation: components.NewButtonAnimationDefault(),
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

func (r *RingSizeSelector) Changed() (bool, int) {
	return r.changed, r.Value
}

func (r *RingSizeSelector) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	r.changed = false

	if r.buttonSelect.Clicked() {
		r.selectModal.Modal.SetVisible(true)
	}

	selected, key := r.selectModal.Selected()
	if selected {
		v, _ := strconv.Atoi(key)
		r.Value = v
		r.changed = true
		r.selectModal.Modal.SetVisible(false)
	}

	r.buttonSelect.Text = fmt.Sprintf("Ring size: %s", fmt.Sprint(r.Value))
	r.buttonSelect.Style.Colors = theme.Current.ButtonPrimaryColors
	return r.buttonSelect.Layout(gtx, th)
}
