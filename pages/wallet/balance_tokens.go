package page_wallet

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/io/clipboard"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/containers/bottom_bar"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageBalanceTokens struct {
	isActive   bool
	firstEnter bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	alertBox       *AlertBox
	displayBalance *DisplayBalance
	tokenBar       *TokenBar
	tokenItems     []*TokenListItem
	buttonSettings *components.Button
	buttonRegister *components.Button
	buttonCopyAddr *components.Button

	list *widget.List
}

var _ router.Page = &PageBalanceTokens{}

func NewPageBalanceTokens() *PageBalanceTokens {
	th := app_instance.Theme

	img, err := assets.GetImage("dero.jpg")
	if err != nil {
		log.Fatal(err)
	}

	tokenItems := []*TokenListItem{}
	for i := 0; i < 10; i++ {
		tokenItems = append(tokenItems, &TokenListItem{
			tokenImageItem: NewTokenImageItem(img),
			tokenName:      fmt.Sprintf("Dero %d", i),
			tokenId:        "00000...00000",
			tokenBalance:   "342.35546",
			Clickable:      new(widget.Clickable),
			//ImageClickable: new(widget.Clickable),
		})
	}

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	settingsIcon, _ := widget.NewIcon(icons.ActionSettings)
	buttonSettings := components.NewButton(components.ButtonStyle{
		Icon:      settingsIcon,
		TextColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		Animation: components.NewButtonAnimationScale(.98),
	})

	buttonRegister := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Text:            lang.Translate("REGISTER WALLET"),
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonRegister.Label.Alignment = text.Middle
	buttonRegister.Style.Font.Weight = font.Bold

	textColor := color.NRGBA{R: 0, G: 0, B: 0, A: 100}
	textHoverColor := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

	copyIcon, _ := widget.NewIcon(icons.ContentContentCopy)
	buttonCopyAddr := components.NewButton(components.ButtonStyle{
		Icon:           copyIcon,
		TextColor:      textColor,
		HoverTextColor: &textHoverColor,
	})

	return &PageBalanceTokens{
		displayBalance: NewDisplayBalance(th),
		tokenBar:       NewTokenBar(th),
		alertBox:       NewAlertBox(),
		tokenItems:     tokenItems,
		firstEnter:     true,
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		list:           list,
		buttonSettings: buttonSettings,
		buttonRegister: buttonRegister,
		buttonCopyAddr: buttonCopyAddr,
	}
}

func (p *PageBalanceTokens) IsActive() bool {
	return p.isActive
}

func (p *PageBalanceTokens) Enter() {
	p.isActive = true

	if !p.firstEnter {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}

	p.ResetWalletHeader()
	page_instance.header.ButtonRight = p.buttonSettings
	bottom_bar.Instance.SetButtonActive(bottom_bar.BUTTON_WALLET)
	p.firstEnter = false
}

func (p *PageBalanceTokens) ResetWalletHeader() {
	openedWallet := wallet_manager.OpenedWallet
	page_instance.header.Title = fmt.Sprintf("%s [%s]", lang.Translate("Wallet"), openedWallet.Info.Name)

	th := app_instance.Theme
	page_instance.header.ButtonRight = nil
	page_instance.header.Subtitle = func(gtx layout.Context) layout.Dimensions {
		walletAddr := openedWallet.Info.Addr
		if p.buttonCopyAddr.Clickable.Clicked() {
			clipboard.WriteOp{
				Text: walletAddr,
			}.Add(gtx.Ops)
			notification_modals.InfoInstance.SetText(lang.Translate("Clipboard"), lang.Translate("Addr copied to clipboard"))
			notification_modals.InfoInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		}

		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				walletAddr := utils.ReduceAddr(walletAddr)
				label := material.Label(th, unit.Sp(16), walletAddr)
				label.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 200}
				return label.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(5)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X = gtx.Dp(18)
				gtx.Constraints.Max.Y = gtx.Dp(18)
				return p.buttonCopyAddr.Layout(gtx, th)
			}),
		)
	}
}

func (p *PageBalanceTokens) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
	page_instance.header.ButtonRight = nil
}

func (p *PageBalanceTokens) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	{
		state := p.animationEnter.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	{
		state := p.animationLeave.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}

		if state.Finished {
			p.isActive = false
			op.InvalidateOp{}.Add(gtx.Ops)
		}
	}

	if p.buttonSettings.Clickable.Clicked() {
		page_instance.pageRouter.SetCurrent(PAGE_SETTINGS)
	}

	var wallet *walletapi.Wallet_Memory
	if wallet_manager.OpenedWallet != nil {
		wallet = wallet_manager.OpenedWallet.Memory
	}

	if walletapi.Connected && wallet != nil {
		isRegistered := wallet.IsRegistered()

		if !isRegistered {
			p.alertBox.SetText("This wallet is not registered on the blockchain.")
			p.alertBox.SetVisible(true)
		} else {
			p.alertBox.SetVisible(false)
		}
	} else {
		p.alertBox.SetText("Wallet is not connected to a node.")
		p.alertBox.SetVisible(true)
	}

	if p.buttonRegister.Clickable.Clicked() {
		page_instance.pageRouter.SetCurrent(PAGE_REGISTER_WALLET)
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.alertBox.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			if walletapi.Connected && wallet != nil {
				isRegistered := wallet.IsRegistered()

				if !isRegistered {
					return layout.Inset{
						Top: unit.Dp(0), Bottom: unit.Dp(20),
						Left: unit.Dp(30), Right: unit.Dp(30),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return p.buttonRegister.Layout(gtx, th)
					})
				}
			}

			return layout.Dimensions{}
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.displayBalance.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.tokenBar.Layout(gtx, th, p.tokenItems)
		},
	}

	for _, item := range p.tokenItems {
		widgets = append(widgets, item.Layout)

		if item.Clickable.Clicked() {
			page_instance.pageRouter.SetCurrent(PAGE_SC_TOKEN)
		}
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
	})

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return widgets[index](gtx)
	})
}

type AlertBox struct {
	iconWarning *widget.Icon
	visible     bool
	text        string
}

func NewAlertBox() *AlertBox {
	iconWarning, _ := widget.NewIcon(icons.AlertWarning)
	return &AlertBox{
		iconWarning: iconWarning,
	}
}

func (n *AlertBox) SetText(value string) {
	n.text = value
}

func (n *AlertBox) SetVisible(visible bool) {
	n.visible = visible
}

func (n *AlertBox) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if !n.visible {
		return layout.Dimensions{}
	}

	border := widget.Border{Color: color.NRGBA{A: 100}, CornerRadius: unit.Dp(5), Width: unit.Dp(1)}

	return layout.Inset{
		Top: unit.Dp(0), Bottom: unit.Dp(20),
		Left: unit.Dp(30), Right: unit.Dp(30),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return n.iconWarning.Layout(gtx, color.NRGBA{A: 100})
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						label := material.Label(th, unit.Sp(14), n.text)
						label.Color = color.NRGBA{A: 200}
						return label.Layout(gtx)
					}),
				)
			})
		})
	})
}

type DisplayBalance struct {
	buttonSend        *components.Button
	buttonReceive     *components.Button
	buttonHideBalance *components.Button

	hideBalanceIcon *widget.Icon
	showBalanceIcon *widget.Icon
}

func NewDisplayBalance(th *material.Theme) *DisplayBalance {
	sendIcon, _ := widget.NewIcon(icons.NavigationArrowUpward)
	buttonSend := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Text:            lang.Translate("SEND"),
		Icon:            sendIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonSend.Label.Alignment = text.Middle
	buttonSend.Style.Font.Weight = font.Bold

	receiveIcon, _ := widget.NewIcon(icons.NavigationArrowDownward)
	buttonReceive := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Text:            lang.Translate("RECEIVE"),
		Icon:            receiveIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonReceive.Label.Alignment = text.Middle
	buttonReceive.Style.Font.Weight = font.Bold

	hideBalanceIcon, _ := widget.NewIcon(icons.ActionVisibility)
	showBalanceIcon, _ := widget.NewIcon(icons.ActionVisibilityOff)
	buttonHideBalance := components.NewButton(components.ButtonStyle{
		Icon:      hideBalanceIcon,
		TextColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
	})

	return &DisplayBalance{
		buttonSend:        buttonSend,
		buttonReceive:     buttonReceive,
		buttonHideBalance: buttonHideBalance,
		hideBalanceIcon:   hideBalanceIcon,
		showBalanceIcon:   showBalanceIcon,
	}
}

func (d *DisplayBalance) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Inset{
		Left: unit.Dp(30), Right: unit.Dp(30),
		Top: unit.Dp(0), Bottom: unit.Dp(40),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lblTitle := material.Label(th, unit.Sp(14), lang.Translate("Available Balance"))
				lblTitle.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 150}

				return lblTitle.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						amount := utils.ShiftNumber{Number: 100000, Decimals: 5}.Format()
						lblAmount := material.Label(th, unit.Sp(34), amount)
						lblAmount.Font.Weight = font.Bold
						dims := lblAmount.Layout(gtx)

						if settings.App.HideBalance {
							paint.FillShape(gtx.Ops, color.NRGBA{R: 200, G: 200, B: 200, A: 255}, clip.Rect{
								Max: dims.Size,
							}.Op())
						}

						return dims
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.Y = gtx.Dp(50)
						gtx.Constraints.Max.X = gtx.Dp(50)

						if settings.App.HideBalance {
							d.buttonHideBalance.Style.Icon = d.showBalanceIcon
						} else {
							d.buttonHideBalance.Style.Icon = d.hideBalanceIcon
						}

						if d.buttonHideBalance.Clickable.Clicked() {
							settings.App.HideBalance = !settings.App.HideBalance
							settings.Save()
							op.InvalidateOp{}.Add(gtx.Ops)
						}

						return d.buttonHideBalance.Layout(gtx, th)
					}),
				)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.Y = gtx.Dp(40)
						return d.buttonSend.Layout(gtx, th)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(15)}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.Y = gtx.Dp(40)

						return d.buttonReceive.Layout(gtx, th)
					}),
				)
			}),
		)
	})
}

type TokenBar struct {
	buttonAddToken  *components.Button
	buttonListToken *components.Button
}

func NewTokenBar(th *material.Theme) *TokenBar {
	addIcon, _ := widget.NewIcon(icons.ContentAddBox)
	buttonAddToken := components.NewButton(components.ButtonStyle{
		Icon:           addIcon,
		TextColor:      color.NRGBA{A: 100},
		HoverTextColor: &color.NRGBA{A: 255},
		Animation:      components.NewButtonAnimationScale(.92),
	})

	listIcon, _ := widget.NewIcon(icons.ActionViewList)
	buttonListToken := components.NewButton(components.ButtonStyle{
		Icon:           listIcon,
		TextColor:      color.NRGBA{A: 100},
		HoverTextColor: &color.NRGBA{A: 255},
		Animation:      components.NewButtonAnimationScale(.92),
	})

	return &TokenBar{
		buttonAddToken:  buttonAddToken,
		buttonListToken: buttonListToken,
	}
}

func (t *TokenBar) Layout(gtx layout.Context, th *material.Theme, items []*TokenListItem) layout.Dimensions {
	cl := clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Dp(1))}.Push(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{R: 0, G: 0, B: 0, A: 50}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl.Pop()

	return layout.Inset{
		Left: unit.Dp(30), Right: unit.Dp(30),
		Top: unit.Dp(30), Bottom: unit.Dp(20),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						labelTokens := material.Label(th, unit.Sp(17), lang.Translate("YOUR TOKENS"))
						labelTokens.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 200}
						labelTokens.Font.Weight = font.Bold
						return labelTokens.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(35)
						gtx.Constraints.Min.Y = gtx.Dp(35)
						return t.buttonAddToken.Layout(gtx, th)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(35)
						gtx.Constraints.Min.Y = gtx.Dp(35)
						return t.buttonListToken.Layout(gtx, th)
					}),
				)
			}),
		)
	})
}

type TokenImageItem struct {
	Image     *components.Image
	Clickable *widget.Clickable

	AnimationEnter   *animation.Animation
	AnimationLeave   *animation.Animation
	hoverSwitchState bool
}

func NewTokenImageItem(src image.Image) *TokenImageItem {
	image := &components.Image{
		Src: paint.NewImageOp(src),
		Fit: components.Cover,
		RNW: unit.Dp(10),
		RNE: unit.Dp(10),
		RSW: unit.Dp(10),
		RSE: unit.Dp(10),
	}

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 1.1, .1, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1.1, 1, .1, ease.Linear),
	))

	return &TokenImageItem{
		Image:          image,
		Clickable:      new(widget.Clickable),
		AnimationEnter: animationEnter,
		AnimationLeave: animationLeave,
	}
}

func (item *TokenImageItem) Layout(gtx layout.Context) layout.Dimensions {
	return item.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		{
			state := item.AnimationEnter.Update(gtx)
			if state.Active {
				item.Image.Transform = func(dims layout.Dimensions, trans f32.Affine2D) f32.Affine2D {
					pt := dims.Size.Div(2)
					origin := f32.Pt(float32(pt.X), float32(pt.Y))
					return trans.Scale(origin, f32.Pt(state.Value, state.Value))
				}
			}
		}

		{
			state := item.AnimationLeave.Update(gtx)
			if state.Active {
				item.Image.Transform = func(dims layout.Dimensions, trans f32.Affine2D) f32.Affine2D {
					pt := dims.Size.Div(2)
					origin := f32.Pt(float32(pt.X), float32(pt.Y))
					return trans.Scale(origin, f32.Pt(state.Value, state.Value))
				}
			}
		}

		/*
			v1, a1, _ := item.AnimationEnter.Update(gtx)
			v2, a2, _ := item.AnimationLeave.Update(gtx)

			item.Image.Transform = func(dims layout.Dimensions, trans f32.Affine2D) f32.Affine2D {
				pt := dims.Size.Div(2)
				origin := f32.Pt(float32(pt.X), float32(pt.Y))

				if a1 {
					trans = trans.Scale(origin, f32.Pt(v1, v1))
				}

				if a2 {
					trans = trans.Scale(origin, f32.Pt(v2, v2))
				}

				return trans
			}*/

		if item.Clickable.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
		}

		if item.Clickable.Hovered() && !item.hoverSwitchState {
			item.hoverSwitchState = true
			item.AnimationEnter.Start()
			item.AnimationLeave.Reset()
		}

		if !item.Clickable.Hovered() && item.hoverSwitchState {
			item.hoverSwitchState = false
			item.AnimationLeave.Start()
			item.AnimationEnter.Reset()
		}

		return item.Image.Layout(gtx)
	})
}

type TokenListItem struct {
	tokenName    string
	tokenId      string
	tokenBalance string

	Clickable      *widget.Clickable
	tokenImageItem *TokenImageItem
}

func (item *TokenListItem) Layout(gtx layout.Context) layout.Dimensions {
	if item.Clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
	}

	return layout.Inset{
		Top: unit.Dp(0), Bottom: unit.Dp(10),
		Right: unit.Dp(30), Left: unit.Dp(30),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		th := app_instance.Theme
		m := op.Record(gtx.Ops)
		dims := item.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top: unit.Dp(10), Bottom: unit.Dp(10),
				Left: unit.Dp(15), Right: unit.Dp(15),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.X = gtx.Dp(50)
						gtx.Constraints.Max.Y = gtx.Dp(50)
						return item.tokenImageItem.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis:      layout.Horizontal,
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										label := material.Label(th, unit.Sp(18), item.tokenName)
										label.Font.Weight = font.Bold
										return label.Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										label := material.Label(th, unit.Sp(14), item.tokenId)
										label.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 150}
										return label.Layout(gtx)
									}),
								)
							}),
							layout.Flexed(1, layout.Spacer{}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								label := material.Label(th, unit.Sp(18), item.tokenBalance)
								label.Font.Weight = font.Bold
								return label.Layout(gtx)
							}),
						)
					}),
				)
			})
		})
		c := m.Stop()

		paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			clip.RRect{
				Rect: image.Rectangle{Max: dims.Size},
				SE:   gtx.Dp(10),
				NW:   gtx.Dp(10),
				NE:   gtx.Dp(10),
				SW:   gtx.Dp(10),
			}.Op(gtx.Ops))

		c.Add(gtx.Ops)

		return dims
	})
	/*
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return dims
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
		)


			if item.Clickable.Hovered() {
				pointer.CursorPointer.Add(gtx.Ops)
				bounds := image.Rect(0, 0, dims.Size.X, dims.Size.Y)
				paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 100},
					clip.UniformRRect(bounds, 10).Op(gtx.Ops),
				)
			}*/
}
