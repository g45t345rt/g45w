package page_wallet

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageSCViewCode struct {
	isActive bool

	headerPageAnimation *prefabs.PageHeaderAnimation
	codeEditor          *widget.Editor
	list                *widget.List
}

var _ router.Page = &PageSCViewCode{}

func NewPageSCViewCode() *PageSCViewCode {

	list := new(widget.List)
	list.Axis = layout.Vertical

	validIcon, _ := widget.NewIcon(icons.ActionCheckCircle)
	buttonExecute := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      validIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonExecute.Label.Alignment = text.Middle
	buttonExecute.Style.Font.Weight = font.Bold

	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_SC_VIEW_CODE)
	return &PageSCViewCode{
		headerPageAnimation: headerPageAnimation,
		list:                list,
	}
}

func (p *PageSCViewCode) IsActive() bool {
	return p.isActive
}

func (p *PageSCViewCode) SetCode(code string) {
	// instanciating everytime - it crash when I call SetText with text already defined
	p.codeEditor = &widget.Editor{ReadOnly: true}
	p.codeEditor.SetText(code)
}

func (p *PageSCViewCode) Enter() {
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)

	page_instance.header.Title = func() string {
		return lang.Translate("Code")
	}

	page_instance.header.LeftLayout = nil
	page_instance.header.RightLayout = nil
}

func (p *PageSCViewCode) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)
}

func (p *PageSCViewCode) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	return layout.Inset{
		Left: theme.PagePadding, Right: theme.PagePadding,
		Bottom: unit.Dp(30),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		r := op.Record(gtx.Ops)
		dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if p.codeEditor == nil {
				return layout.Dimensions{}
			}

			editor := material.Editor(th, p.codeEditor, "")
			editor.TextSize = unit.Sp(12)
			return editor.Layout(gtx)
		})
		c := r.Stop()

		paint.FillShape(gtx.Ops, theme.Current.BgColor, clip.RRect{
			Rect: image.Rectangle{Max: dims.Size},
			NW:   gtx.Dp(10), NE: gtx.Dp(10),
			SE: gtx.Dp(10), SW: gtx.Dp(10),
		}.Op(gtx.Ops))

		c.Add(gtx.Ops)

		return dims
	})
}
