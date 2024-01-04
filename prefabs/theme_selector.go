package prefabs

import (
	"fmt"
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
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

type ThemeSelector struct {
	buttonSelect *components.Button
	items        []*listselect_modal.SelectListItem

	Changed bool
	Key     string
}

func NewThemeSelector(defaultThemeKey string) *ThemeSelector {
	colorIcon, _ := widget.NewIcon(icons.EditorFormatColorFill)
	buttonSelect := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		TextSize:  unit.Sp(16),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Icon:      colorIcon,
		IconGap:   unit.Dp(10),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonSelect.Label.Alignment = text.Middle
	buttonSelect.Style.Font.Weight = font.Bold

	items := []*listselect_modal.SelectListItem{}

	for _, theme := range theme.Themes {
		indicatorColor := theme.IndicatorColor
		name := theme.Name
		items = append(items, listselect_modal.NewSelectListItem(theme.Key, func(gtx layout.Context, th *material.Theme) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					rect := image.Rectangle{Max: image.Pt(gtx.Dp(25), gtx.Dp(25))}
					paint.FillShape(gtx.Ops, indicatorColor, clip.UniformRRect(rect, gtx.Dp(5)).Op(gtx.Ops))

					return layout.Dimensions{Size: rect.Max}
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(20), lang.Translate(name))
					return lbl.Layout(gtx)
				}),
			)
		}))
	}

	defaultTheme := lang.Translate(defaultThemeKey)
	return &ThemeSelector{
		buttonSelect: buttonSelect,
		items:        items,
		Key:          defaultTheme,
	}
}

func (t *ThemeSelector) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	t.Changed = false

	if t.buttonSelect.Clicked(gtx) {
		go func() {
			keyChan := listselect_modal.Instance.Open(t.items)

			for key := range keyChan {
				t.Changed = true
				t.Key = key
			}
		}()
	}

	text := ""
	for _, theme := range theme.Themes {
		if theme.Key == t.Key {
			text = lang.Translate(theme.Name)
		}
	}

	t.buttonSelect.Text = fmt.Sprintf("%s: %s", lang.Translate("Theme"), text)
	t.buttonSelect.Style.Colors = theme.Current.ButtonPrimaryColors
	return t.buttonSelect.Layout(gtx, th)
}
