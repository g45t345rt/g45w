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
	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/listselect_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type LangSelector struct {
	buttonSelect *components.Button
	items        []*listselect_modal.SelectListItem

	Key     string
	Changed bool
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

	var items []*listselect_modal.SelectListItem

	for _, language := range lang.Languages {
		img, _ := assets.GetImage(language.ImgPath)

		langImg := &components.Image{
			Src:      paint.NewImageOp(img),
			Fit:      components.Cover,
			Position: layout.Center,
			Rounded:  components.UniformRounded(unit.Dp(5)),
		}

		name := language.Name
		items = append(items, listselect_modal.NewSelectListItem(language.Key, func(gtx layout.Context, th *material.Theme) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = gtx.Dp(45)
					gtx.Constraints.Max.Y = gtx.Dp(30)
					return langImg.Layout(gtx, nil)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(20), lang.Translate(name))
					return lbl.Layout(gtx)
				}),
			)
		}))
	}

	return &LangSelector{
		buttonSelect: buttonSelect,
		items:        items,
		Key:          defaultLangKey,
	}
}

func (l *LangSelector) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	l.Changed = false
	if l.buttonSelect.Clicked(gtx) {
		go func() {
			keyChan := listselect_modal.Instance.Open(l.items)
			for key := range keyChan {
				l.Key = key
				l.Changed = true
			}
		}()
	}

	text := ""
	for _, language := range lang.Languages {
		if language.Key == l.Key {
			text = lang.Translate(language.Name)
		}
	}

	l.buttonSelect.Text = fmt.Sprintf("%s: %s", lang.Translate("Language"), text)
	l.buttonSelect.Style.Colors = theme.Current.ButtonPrimaryColors
	return l.buttonSelect.Layout(gtx, th)
}
