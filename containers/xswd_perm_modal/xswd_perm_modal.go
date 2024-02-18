package xswd_perm_modal

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/creachadair/jrpc2"
	"github.com/deroproject/derohe/walletapi/xswd"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type XSWDPermModal struct {
	Modal *components.Modal

	req               *jrpc2.Request
	app               *xswd.ApplicationData
	buttonAllow       *components.Button
	buttonDeny        *components.Button
	buttonAlwaysAllow *components.Button
	buttonAlwaysDeny  *components.Button
	buttonClose       *components.Button

	permChan chan xswd.Permission
}

var Instance *XSWDPermModal

func LoadInstance() {
	modal := components.NewModal(components.ModalStyle{
		CloseOnOutsideClick: false,
		CloseOnInsideClick:  false,
		Direction:           layout.N,
		Rounded:             components.UniformRounded(unit.Dp(10)),
		Inset:               layout.UniformInset(theme.PagePadding),
		Animation:           components.NewModalAnimationDown(),
	})

	allowIcon, _ := widget.NewIcon(icons.ActionThumbUp)
	buttonAllow := components.NewButton(components.ButtonStyle{
		Icon:    allowIcon,
		IconGap: unit.Dp(10),
		Rounded: components.UniformRounded(unit.Dp(5)),
		Border: widget.Border{
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
		TextSize:  unit.Sp(14),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonAllow.Label.Alignment = text.Middle
	buttonAllow.Style.Font.Weight = font.Bold

	denyIcon, _ := widget.NewIcon(icons.ActionThumbDown)
	buttonDeny := components.NewButton(components.ButtonStyle{
		Icon:     denyIcon,
		IconGap:  unit.Dp(10),
		Rounded:  components.UniformRounded(unit.Dp(5)),
		TextSize: unit.Sp(14),
		Border: widget.Border{
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonDeny.Label.Alignment = text.Middle
	buttonDeny.Style.Font.Weight = font.Bold

	buttonAlwaysAllow := components.NewButton(components.ButtonStyle{
		Icon:     allowIcon,
		IconGap:  unit.Dp(10),
		Rounded:  components.UniformRounded(unit.Dp(5)),
		TextSize: unit.Sp(14),
		Border: widget.Border{
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonAlwaysAllow.Label.Alignment = text.Middle
	buttonAlwaysAllow.Style.Font.Weight = font.Bold

	buttonAlwaysDeny := components.NewButton(components.ButtonStyle{
		Icon:     denyIcon,
		IconGap:  unit.Dp(10),
		Rounded:  components.UniformRounded(unit.Dp(5)),
		TextSize: unit.Sp(14),
		Border: widget.Border{
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonAlwaysDeny.Label.Alignment = text.Middle
	buttonAlwaysDeny.Style.Font.Weight = font.Bold

	closeIcon, _ := widget.NewIcon(icons.NavigationCancel)
	buttonClose := components.NewButton(components.ButtonStyle{
		Icon:      closeIcon,
		Animation: components.NewButtonAnimationDefault(),
	})

	Instance = &XSWDPermModal{
		Modal:             modal,
		buttonAllow:       buttonAllow,
		buttonDeny:        buttonDeny,
		buttonAlwaysAllow: buttonAlwaysAllow,
		buttonAlwaysDeny:  buttonAlwaysDeny,
		buttonClose:       buttonClose,
	}

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 3,
		Layout: func(gtx layout.Context, th *material.Theme) {
			Instance.Layout(gtx, th)
		},
	})
}

func (c *XSWDPermModal) Open(app *xswd.ApplicationData, req *jrpc2.Request) chan xswd.Permission {
	c.app = app
	c.req = req

	c.Modal.SetVisible(true)
	c.permChan = make(chan xswd.Permission)
	return c.permChan
}

func (c *XSWDPermModal) set(perm xswd.Permission) {
	c.permChan <- perm
	c.Modal.SetVisible(false)
	close(c.permChan)
	app_instance.Window.Invalidate()
}

func (c *XSWDPermModal) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if c.buttonAllow.Clicked(gtx) {
		go c.set(xswd.Allow)
	}

	if c.buttonDeny.Clicked(gtx) {
		go c.set(xswd.Deny)
	}

	if c.buttonAlwaysAllow.Clicked(gtx) {
		go c.set(xswd.AlwaysAllow)
	}

	if c.buttonAlwaysDeny.Clicked(gtx) {
		go c.set(xswd.AlwaysDeny)
	}

	if c.buttonClose.Clicked(gtx) {
		go c.set(xswd.Deny)
	}

	c.Modal.Style.Colors = theme.Current.ModalColors
	return c.Modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									lbl := material.Label(th, unit.Sp(22), c.app.Name)
									lbl.Font.Weight = font.Bold
									return lbl.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									c.buttonClose.Style.Colors = theme.Current.ModalButtonColors
									return c.buttonClose.Layout(gtx, th)
								}),
							)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									txt := lang.Translate("From:")
									lbl := material.Label(th, unit.Sp(16), txt)
									lbl.Color = theme.Current.TextMuteColor
									return lbl.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Width: unit.Dp(3)}.Layout),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									lbl := material.Label(th, unit.Sp(16), c.app.Url)
									return lbl.Layout(gtx)
								}),
							)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									txt := lang.Translate("Method:")
									lbl := material.Label(th, unit.Sp(16), txt)
									lbl.Color = theme.Current.TextMuteColor
									return lbl.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Width: unit.Dp(3)}.Layout),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									lbl := material.Label(th, unit.Sp(16), c.req.Method())
									return lbl.Layout(gtx)
								}),
							)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									txt := lang.Translate("Parameters:")
									lbl := material.Label(th, unit.Sp(16), txt)
									lbl.Color = theme.Current.TextMuteColor
									return lbl.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Width: unit.Dp(3)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									params := c.req.ParamString()
									lbl := material.Label(th, unit.Sp(14), params)
									return lbl.Layout(gtx)
								}),
							)
						}),
					)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					c.buttonAllow.Style.Colors = theme.Current.ButtonSecondaryColors
					c.buttonAllow.Text = lang.Translate("Allow")
					return c.buttonAllow.Layout(gtx, th)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					c.buttonAlwaysAllow.Style.Colors = theme.Current.ButtonSecondaryColors
					c.buttonAlwaysAllow.Text = lang.Translate("Allow Always")
					return c.buttonAlwaysAllow.Layout(gtx, th)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return prefabs.Divider(gtx, unit.Dp(5))
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					c.buttonDeny.Style.Colors = theme.Current.ButtonSecondaryColors
					c.buttonDeny.Text = lang.Translate("Deny")
					return c.buttonDeny.Layout(gtx, th)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					c.buttonAlwaysDeny.Style.Colors = theme.Current.ButtonSecondaryColors
					c.buttonAlwaysDeny.Text = lang.Translate("Deny Always")
					return c.buttonAlwaysDeny.Layout(gtx, th)
				}),
			)
		})
	})
}
