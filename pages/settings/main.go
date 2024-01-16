package page_settings

import (
	"image/color"
	"os"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/android_background_service"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_icons"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/confirm_modal"
	"github.com/g45t345rt/g45w/containers/notification_modal"
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

	langSelector                   *prefabs.LangSelector
	themeSelector                  *prefabs.ThemeSelector
	buttonInfo                     *components.Button
	buttonIpfsGateway              *components.Button
	buttonDonation                 *components.Button
	androidBackgroundServiceSwitch *AndroidBackgroundServiceSwitch
	testnetSwitch                  *TestnetSwitch
}

var _ router.Page = &PageMain{}

func NewPageFront() *PageMain {
	defaultLangKey := settings.App.Language
	defaultThemeKey := settings.App.Theme
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

	donationIcon, _ := widget.NewIcon(app_icons.Donation)
	buttonDonation := components.NewButton(components.ButtonStyle{
		Icon:      donationIcon,
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
	buttonDonation.Label.Alignment = text.Middle
	buttonDonation.Style.Font.Weight = font.Bold

	list := new(widget.List)
	list.Axis = layout.Vertical

	return &PageMain{
		list:           list,
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		langSelector:                   langSelector,
		themeSelector:                  themeSelector,
		buttonInfo:                     buttonInfo,
		buttonIpfsGateway:              buttonIpfsGateway,
		buttonDonation:                 buttonDonation,
		androidBackgroundServiceSwitch: NewAndroidBackgroundServiceSwitch(),
		testnetSwitch:                  NewTestnetSwitch(),
	}
}

func (p *PageMain) IsActive() bool {
	return p.isActive
}

func (p *PageMain) Enter() {
	p.isActive = true
	page_instance.header.Title = func() string { return lang.Translate("Settings") }
	page_instance.header.Subtitle = nil
	page_instance.header.LeftLayout = nil

	if !page_instance.header.IsHistory(PAGE_MAIN) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}

	p.testnetSwitch.Load()
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

	if p.buttonInfo.Clicked(gtx) {
		page_instance.pageRouter.SetCurrent(PAGE_APP_INFO)
		page_instance.header.AddHistory(PAGE_APP_INFO)
	}

	if p.buttonIpfsGateway.Clicked(gtx) {
		page_instance.pageRouter.SetCurrent(PAGE_IPFS_GATEWAYS)
		page_instance.header.AddHistory(PAGE_IPFS_GATEWAYS)
	}

	if p.buttonDonation.Clicked(gtx) {
		page_instance.pageRouter.SetCurrent(PAGE_DONATION)
		page_instance.header.AddHistory(PAGE_DONATION)
	}

	if p.langSelector.Changed {
		go func() {
			settings.App.Language = p.langSelector.Key
			err := settings.Save()
			if err != nil {
				notification_modal.Open(notification_modal.Params{
					Type:  notification_modal.ERROR,
					Title: lang.Translate("Error"),
					Text:  err.Error(),
				})
			} else {
				lang.Current = settings.App.Language
			}
		}()
	}

	if p.themeSelector.Changed {
		go func() {
			settings.App.Theme = p.themeSelector.Key
			err := settings.Save()
			if err != nil {
				notification_modal.Open(notification_modal.Params{
					Type:  notification_modal.ERROR,
					Title: lang.Translate("Error"),
					Text:  err.Error(),
				})
			} else {
				theme.Current = theme.Get(settings.App.Theme)
			}
		}()
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

	if !settings.App.Testnet {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			p.buttonDonation.Text = lang.Translate("Donate")
			p.buttonDonation.Style.Colors = theme.Current.ButtonSecondaryColors
			return p.buttonDonation.Layout(gtx, th)
		})
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.testnetSwitch.Layout(gtx, th)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.androidBackgroundServiceSwitch.Layout(gtx, th)
	})

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

type AndroidBackgroundServiceSwitch struct {
	switchForeground *widget.Bool
	switchValue      bool
}

func NewAndroidBackgroundServiceSwitch() *AndroidBackgroundServiceSwitch {
	return &AndroidBackgroundServiceSwitch{
		switchForeground: new(widget.Bool),
	}
}

func (a *AndroidBackgroundServiceSwitch) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	available := android_background_service.IsAvailable()
	running, _ := android_background_service.IsRunning()

	if available && running {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(20), lang.Translate("Android background service"))
						lbl.Font.Weight = font.Bold
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						s := material.Switch(th, a.switchForeground, "")
						s.Color = theme.Current.SwitchColors

						// switch does not have a released func
						if a.switchForeground.Value != a.switchValue {
							a.switchValue = a.switchForeground.Value

							go func() {
								if a.switchValue {
									android_background_service.StartForeground()
								} else {
									android_background_service.StopForeground()
								}
							}()
						}

						return s.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(14), lang.Translate("This is required to either interact with XSWD or keep the node connection open."))
						lbl.Color = theme.Current.TextMuteColor
						return lbl.Layout(gtx)
					}),
				)
			}),
		)
	}

	return layout.Dimensions{}
}

type TestnetSwitch struct {
	switchTestnet *widget.Bool
	switchValue   bool
}

func NewTestnetSwitch() *TestnetSwitch {
	return &TestnetSwitch{
		switchTestnet: new(widget.Bool),
	}
}

func (a *TestnetSwitch) Load() {
	a.switchTestnet.Value = settings.App.Testnet
	a.switchValue = settings.App.Testnet
}

func (a *TestnetSwitch) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(20), lang.Translate("Testnet"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					s := material.Switch(th, a.switchTestnet, "")
					s.Color = theme.Current.SwitchColors

					if a.switchTestnet.Value != a.switchValue {
						a.switchValue = a.switchTestnet.Value

						go func() {
							yes := <-confirm_modal.Instance.Open(confirm_modal.ConfirmText{
								Title:  lang.Translate("Are you sure?"),
								Prompt: lang.Translate("The application will close, and you'll need to restart it."),
							})
							if yes {
								settings.App.Testnet = a.switchValue
								settings.Save()
								os.Exit(0)
							} else {
								a.switchTestnet.Value = !a.switchTestnet.Value
								a.switchValue = !a.switchValue
							}
						}()
					}

					return s.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), lang.Translate("Switching the application to Testnet mode enables connection with Testnet nodes and creates another instance for wallets."))
					lbl.Color = theme.Current.TextMuteColor
					return lbl.Layout(gtx)
				}),
			)
		}),
	)
}
