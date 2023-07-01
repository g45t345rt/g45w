package prefabs

import (
	"fmt"
	"image/color"
	"log"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/components"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type LangSelector struct {
	buttonSelect *components.Button
	selectModal  *SelectModal

	changed bool
	Value   string
}

func NewLangSelector(defaultLangKey string) *LangSelector {
	langIcon, _ := widget.NewIcon(icons.ActionLanguage)
	buttonSelect := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(16),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Icon:            langIcon,
		IconGap:         unit.Dp(10),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonSelect.Label.Alignment = text.Middle
	buttonSelect.Style.Font.Weight = font.Bold

	items := []*SelectListItem{}
	th := app_instance.Theme
	w := app_instance.Window

	languages := lang.SupportedLanguages
	for _, language := range languages {
		img, err := assets.GetImage(language.ImgPath)
		if err != nil {
			log.Fatal(err)
		}

		langImg := &components.Image{
			Src:      paint.NewImageOp(img),
			Fit:      components.Cover,
			Position: layout.Center,
			Rounded:  components.UniformRounded(unit.Dp(5)),
		}

		items = append(items, NewSelectListItem(language.Key, func(gtx layout.Context, index int) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = gtx.Dp(45)
					gtx.Constraints.Max.Y = gtx.Dp(30)
					return langImg.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					name := languages[index].Name
					lbl := material.Label(th, unit.Sp(20), lang.Translate(name))
					return lbl.Layout(gtx)
				}),
			)
		}))
	}

	selectModal := NewSelectModal(w)
	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			selectModal.Layout(gtx, th, items)
		},
	})

	defaultLanguage := lang.Translate(defaultLangKey)
	r := &LangSelector{
		buttonSelect: buttonSelect,
		selectModal:  selectModal,
		Value:        defaultLanguage,
	}

	return r
}

func (r *LangSelector) Changed() bool {
	return r.changed
}

func (r *LangSelector) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	r.changed = false

	if r.buttonSelect.Clickable.Clicked() {
		r.selectModal.modal.SetVisible(true)
	}

	selected := r.selectModal.Selected()
	if selected {
		r.Value = r.selectModal.SelectedKey
		r.changed = true
		r.selectModal.modal.SetVisible(false)
	}

	value := r.Value
	for _, language := range lang.SupportedLanguages {
		if language.Key == r.Value {
			value = lang.Translate(language.Name)
		}
	}

	r.buttonSelect.Text = fmt.Sprintf("%s: %s", lang.Translate("Language"), value)
	return r.buttonSelect.Layout(gtx, th)
}
