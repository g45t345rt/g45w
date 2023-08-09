package prefabs

import (
	"fmt"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type LangSelector struct {
	ButtonSelect *components.Button
	SelectModal  *SelectModal

	changed bool
	Value   string
}

func NewLangSelector(defaultLangKey string) *LangSelector {
	langIcon, _ := widget.NewIcon(icons.ActionLanguage)
	buttonSelect := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		TextSize:  unit.Sp(16),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Icon:      langIcon,
		IconGap:   unit.Dp(10),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonSelect.Label.Alignment = text.Middle
	buttonSelect.Style.Font.Weight = font.Bold

	items := []*SelectListItem{}

	for _, language := range lang.Languages {
		img, _ := assets.GetImage(language.ImgPath)

		langImg := &components.Image{
			Src:      paint.NewImageOp(img),
			Fit:      components.Cover,
			Position: layout.Center,
			Rounded:  components.UniformRounded(unit.Dp(5)),
		}

		name := language.Name
		items = append(items, NewSelectListItem(language.Key, func(gtx layout.Context, index int, th *material.Theme) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = gtx.Dp(45)
					gtx.Constraints.Max.Y = gtx.Dp(30)
					return langImg.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(20), lang.Translate(name))
					return lbl.Layout(gtx)
				}),
			)
		}))
	}

	selectModal := NewSelectModal()
	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			selectModal.Layout(gtx, th, items)
		},
	})

	defaultLanguage := lang.Translate(defaultLangKey)
	r := &LangSelector{
		ButtonSelect: buttonSelect,
		SelectModal:  selectModal,
		Value:        defaultLanguage,
	}

	return r
}

func (r *LangSelector) Changed() bool {
	return r.changed
}

func (r *LangSelector) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	r.changed = false

	if r.ButtonSelect.Clicked() {
		r.SelectModal.Modal.SetVisible(true)
	}

	selected, key := r.SelectModal.Selected()
	if selected {
		r.Value = key
		r.changed = true
		r.SelectModal.Modal.SetVisible(false)
	}

	value := r.Value
	for _, language := range lang.Languages {
		if language.Key == r.Value {
			value = lang.Translate(language.Name)
		}
	}

	r.ButtonSelect.Text = fmt.Sprintf("%s: %s", lang.Translate("Language"), value)
	r.ButtonSelect.Style.Colors = theme.Current.ButtonPrimaryColors
	return r.ButtonSelect.Layout(gtx, th)
}
