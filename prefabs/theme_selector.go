package prefabs

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var themes = map[string]string{
	"light": "Light", //@lang.Translate("Light")
	"dark":  "Dark",  //@lang.Translate("Dark")
	"blue":  "Blue",  //@lang.Translate("Blue")
}

type ThemeSelector struct {
	ButtonSelect *components.Button
	SelectModal  *SelectModal

	changed bool
	Value   string
}

func NewThemeSelector(defaultThemeKey string) *ThemeSelector {
	colorIcon, _ := widget.NewIcon(icons.EditorFormatColorFill)
	buttonSelect := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(16),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Icon:            colorIcon,
		IconGap:         unit.Dp(10),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonSelect.Label.Alignment = text.Middle
	buttonSelect.Style.Font.Weight = font.Bold

	items := []*SelectListItem{}

	for key := range themes {
		themeKey := key
		items = append(items, NewSelectListItem(key, func(gtx layout.Context, index int, th *material.Theme) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					rect := image.Rectangle{Max: image.Pt(gtx.Dp(25), gtx.Dp(25))}
					paint.FillShape(gtx.Ops, color.NRGBA{A: 255}, clip.UniformRRect(rect, gtx.Dp(5)).Op(gtx.Ops))

					return layout.Dimensions{Size: rect.Max}
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(20), lang.Translate(themes[themeKey]))
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

	defaultTheme := lang.Translate(defaultThemeKey)
	r := &ThemeSelector{
		ButtonSelect: buttonSelect,
		SelectModal:  selectModal,
		Value:        defaultTheme,
	}

	return r
}

func (r *ThemeSelector) Changed() bool {
	return r.changed
}

func (r *ThemeSelector) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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
	for key, name := range themes {
		if key == r.Value {
			value = lang.Translate(name)
		}
	}

	r.ButtonSelect.Text = fmt.Sprintf("%s: %s", lang.Translate("Theme"), value)
	return r.ButtonSelect.Layout(gtx, th)
}
