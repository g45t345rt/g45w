package prefabs

import (
	"fmt"
	"strconv"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/listselect_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type RingSizeSelector struct {
	buttonSelect *components.Button
	items        []*listselect_modal.SelectListItem

	Size    int
	Changed bool
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
	items := []*listselect_modal.SelectListItem{}

	for i, size := range sizes {
		idx := i
		items = append(items, listselect_modal.NewSelectListItem(size, func(gtx layout.Context, th *material.Theme) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(20), sizes[idx])
			return lbl.Layout(gtx)
		}))
	}

	r := &RingSizeSelector{
		buttonSelect: buttonSelect,
		Size:         defaultSize,
		items:        items,
	}

	return r
}

func (r *RingSizeSelector) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	r.Changed = false

	if r.buttonSelect.Clicked(gtx) {
		go func() {
			keyChan := listselect_modal.Instance.Open(r.items)

			for key := range keyChan {
				size, _ := strconv.Atoi(key)
				r.Size = size
				r.Changed = true
			}
		}()
	}

	text := lang.Translate("Ring size: {}")
	text = strings.Replace(text, "{}", fmt.Sprint(r.Size), -1)
	r.buttonSelect.Text = text
	r.buttonSelect.Style.Colors = theme.Current.ButtonPrimaryColors
	return r.buttonSelect.Layout(gtx, th)
}
