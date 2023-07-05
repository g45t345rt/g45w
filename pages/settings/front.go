package page_settings

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageFront struct {
	isActive       bool
	list           *widget.List
	animationEnter *animation.Animation
	animationLeave *animation.Animation
	langSelector   *prefabs.LangSelector
	buttonInfo     *components.Button
}

var _ router.Page = &PageFront{}

func NewPageFront() *PageFront {
	langKey := "en"
	if settings.App.Language != "" {
		langKey = settings.App.Language
	}

	langSelector := prefabs.NewLangSelector(langKey)

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	infoIcon, _ := widget.NewIcon(icons.ActionInfo)
	buttonInfo := components.NewButton(components.ButtonStyle{
		Icon:            infoIcon,
		TextColor:       color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		BackgroundColor: color.NRGBA{A: 0},
		TextSize:        unit.Sp(16),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
		Border: widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
	})
	buttonInfo.Label.Alignment = text.Middle
	buttonInfo.Style.Font.Weight = font.Bold

	list := new(widget.List)
	list.Axis = layout.Vertical

	return &PageFront{
		list:           list,
		langSelector:   langSelector,
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		buttonInfo:     buttonInfo,
	}
}

func (p *PageFront) IsActive() bool {
	return p.isActive
}

func (p *PageFront) Enter() {
	p.isActive = true
	page_instance.header.SetTitle(lang.Translate("Settings"))

	if !page_instance.header.IsHistory(PAGE_FRONT) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageFront) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *PageFront) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	{
		state := p.animationEnter.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	{
		state := p.animationLeave.Update(gtx)
		if state.Finished {
			p.isActive = false
			op.InvalidateOp{}.Add(gtx.Ops)
		}

		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	if p.buttonInfo.Clickable.Clicked() {
		page_instance.pageRouter.SetCurrent(PAGE_INFO)
		page_instance.header.AddHistory(PAGE_INFO)
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			p.buttonInfo.Text = lang.Translate("App Information")
			return p.buttonInfo.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			if p.langSelector.Changed() {
				settings.App.Language = p.langSelector.Value
				err := settings.Save()
				if err != nil {
					notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
					notification_modals.ErrorInstance.SetVisible(gtx, true, notification_modals.CLOSE_AFTER_DEFAULT)
				} else {
					notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("Language applied."))
					notification_modals.SuccessInstance.SetVisible(gtx, true, notification_modals.CLOSE_AFTER_DEFAULT)
				}
			}

			return p.langSelector.Layout(gtx, th)
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}
