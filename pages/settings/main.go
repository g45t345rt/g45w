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
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/theme"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageMain struct {
	isActive       bool
	list           *widget.List
	animationEnter *animation.Animation
	animationLeave *animation.Animation

	langSelector      *prefabs.LangSelector
	themeSelector     *prefabs.ThemeSelector
	buttonInfo        *components.Button
	buttonIpfsGateway *components.Button
}

var _ router.Page = &PageMain{}

func NewPageFront() *PageMain {
	defaultLangKey := settings.App.LanguageKey
	defaultThemeKey := settings.App.ThemeKey
	langSelector := prefabs.NewLangSelector(defaultLangKey)
	themeSelector := prefabs.NewThemeSelector(defaultThemeKey)

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	infoIcon, _ := widget.NewIcon(icons.ActionInfo)
	buttonInfo := components.NewButton(components.ButtonStyle{
		Icon:      infoIcon,
		TextSize:  unit.Sp(16),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
		Border: widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
	})
	buttonInfo.Label.Alignment = text.Middle
	buttonInfo.Style.Font.Weight = font.Bold

	gatewayIcon, _ := widget.NewIcon(icons.HardwareDeviceHub)
	buttonIpfsGateway := components.NewButton(components.ButtonStyle{
		Icon:      gatewayIcon,
		TextSize:  unit.Sp(16),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
		Border: widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
	})
	buttonIpfsGateway.Label.Alignment = text.Middle
	buttonIpfsGateway.Style.Font.Weight = font.Bold

	list := new(widget.List)
	list.Axis = layout.Vertical

	return &PageMain{
		list:           list,
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		langSelector:      langSelector,
		themeSelector:     themeSelector,
		buttonInfo:        buttonInfo,
		buttonIpfsGateway: buttonIpfsGateway,
	}
}

func (p *PageMain) IsActive() bool {
	return p.isActive
}

func (p *PageMain) Enter() {
	p.isActive = true
	page_instance.header.Title = func() string { return lang.Translate("Settings") }
	page_instance.header.Subtitle = nil
	page_instance.header.ButtonRight = nil

	if !page_instance.header.IsHistory(PAGE_MAIN) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageMain) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *PageMain) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonInfo.Clicked() {
		page_instance.pageRouter.SetCurrent(PAGE_APP_INFO)
		page_instance.header.AddHistory(PAGE_APP_INFO)
	}

	if p.buttonIpfsGateway.Clicked() {
		page_instance.pageRouter.SetCurrent(PAGE_IPFS_GATEWAYS)
		page_instance.header.AddHistory(PAGE_IPFS_GATEWAYS)
	}

	if p.langSelector.Changed() {
		settings.App.LanguageKey = p.langSelector.Value
		err := settings.Save()
		if err != nil {
			notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {
			lang.Current = settings.App.LanguageKey
			notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("Language applied."))
			notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		}
	}

	if p.themeSelector.Changed() {
		settings.App.ThemeKey = p.themeSelector.Value
		err := settings.Save()
		if err != nil {
			notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {
			theme.Current = *theme.Get(settings.App.ThemeKey)
			notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("Theme applied."))
			notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		}
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			p.buttonInfo.Text = lang.Translate("App Information")
			p.buttonInfo.Style.Colors = theme.Current.ButtonSecondaryColors
			return p.buttonInfo.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonIpfsGateway.Text = lang.Translate("IPFS Gateways")
			p.buttonIpfsGateway.Style.Colors = theme.Current.ButtonSecondaryColors
			return p.buttonIpfsGateway.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.langSelector.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.themeSelector.Layout(gtx, th)
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
