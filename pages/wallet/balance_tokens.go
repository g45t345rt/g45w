package page_wallet

import (
	"database/sql"
	"fmt"
	"image"
	"image/color"
	"strings"

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
	crypto "github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/node_status_bar"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/node_manager"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageBalanceTokens struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	alertBox           *AlertBox
	displayBalance     *DisplayBalance
	tokenBar           *TokenBar
	tokenItems         []*TokenListItem
	buttonSettings     *components.Button
	buttonRegister     *components.Button
	buttonCopyAddr     *components.Button
	tabBars            *components.TabBars
	txBar              *TxBar
	txItems            []*TxListItem
	getTransfersParams wallet_manager.GetTransfersParams

	list *widget.List
}

var _ router.Page = &PageBalanceTokens{}

func NewPageBalanceTokens() *PageBalanceTokens {
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

	registerIcon, _ := widget.NewIcon(icons.ActionAssignmentTurnedIn)
	buttonRegister := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
		Icon:            registerIcon,
		IconGap:         unit.Dp(10),
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

	tabBarsItems := []*components.TabBarsItem{
		components.NewTabBarItem("tokens", components.TabBarItemStyle{
			TextSize: unit.Sp(18),
		}),
		components.NewTabBarItem("txs", components.TabBarItemStyle{
			TextSize: unit.Sp(18),
		}),
	}
	defaultTabKey := settings.App.TabBarsKey
	tabBars := components.NewTabBars(defaultTabKey, tabBarsItems)

	txBar := NewTxBar()

	return &PageBalanceTokens{
		displayBalance: NewDisplayBalance(),
		tokenBar:       NewTokenBar(),
		alertBox:       NewAlertBox(),
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		list:           list,
		buttonSettings: buttonSettings,
		buttonRegister: buttonRegister,
		buttonCopyAddr: buttonCopyAddr,
		tabBars:        tabBars,
		txBar:          txBar,
	}
}

func (p *PageBalanceTokens) IsActive() bool {
	return p.isActive
}

func (p *PageBalanceTokens) Enter() {
	p.isActive = true

	if !page_instance.header.IsHistory(PAGE_BALANCE_TOKENS) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}

	p.ResetWalletHeader()
	page_instance.header.ButtonRight = p.buttonSettings
	p.Load()
}

func (p *PageBalanceTokens) Load() error {
	p.LoadTxs()

	err := p.LoadTokens()
	if err != nil {
		return err
	}

	return nil
}

func (p *PageBalanceTokens) LoadTokens() error {
	wallet := wallet_manager.OpenedWallet

	tokens, err := wallet.GetTokens(wallet_manager.GetTokensParams{
		IsFavorite: sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		return err
	}

	tokenItems := []*TokenListItem{}
	imgToken, _ := assets.GetImage("token.png")
	srcImgToken := paint.NewImageOp(imgToken)

	for _, token := range tokens {
		tokenItems = append(tokenItems, NewTokenListItem(token, srcImgToken))
	}

	p.tokenItems = tokenItems
	p.RefreshTokensBalance()

	return nil
}

func (p *PageBalanceTokens) LoadTxs() {
	wallet := wallet_manager.OpenedWallet
	entries := wallet.GetTransfers(crypto.ZEROHASH.String(), p.getTransfersParams)

	txItems := []*TxListItem{}

	imgUp, _ := assets.GetImage("arrow_up_arc.png")
	srcImgUp := paint.NewImageOp(imgUp)

	imgDown, _ := assets.GetImage("arrow_down_arc.png")
	srcImgDown := paint.NewImageOp(imgDown)

	imgCoinbase, _ := assets.GetImage("coinbase.png")
	srcImgCoinbase := paint.NewImageOp(imgCoinbase)

	for _, entry := range entries {
		var img paint.ImageOp

		if entry.Incoming {
			img = srcImgDown
		} else {
			img = srcImgUp
		}

		if entry.Coinbase {
			img = srcImgCoinbase
		}

		txItems = append(txItems, NewTxListItem(entry, img))
	}

	p.txItems = txItems
	p.txBar.txCount = len(entries)
}

func (p *PageBalanceTokens) RefreshTokensBalance() {
	wallet := wallet_manager.OpenedWallet
	for _, tokenItem := range p.tokenItems {
		scId := crypto.HashHexToHash(tokenItem.token.SCID)
		b, _ := wallet.Memory.Get_Balance_scid(scId)
		tokenItem.balance = b
	}
}

func (p *PageBalanceTokens) ResetWalletHeader() {
	openedWallet := wallet_manager.OpenedWallet
	title := fmt.Sprintf("%s [%s]", lang.Translate("Wallet"), openedWallet.Info.Name)
	page_instance.header.SetTitle(title)

	page_instance.header.ButtonRight = nil
	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		walletAddr := openedWallet.Info.Addr
		if p.buttonCopyAddr.Clicked() {
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

	if p.buttonSettings.Clicked() {
		page_instance.pageRouter.SetCurrent(PAGE_SETTINGS)
		page_instance.header.AddHistory(PAGE_SETTINGS)
	}

	if p.buttonRegister.Clicked() {
		page_instance.pageRouter.SetCurrent(PAGE_REGISTER_WALLET)
		page_instance.header.AddHistory(PAGE_REGISTER_WALLET)
	}

	widgets := []layout.Widget{}
	wallet := wallet_manager.OpenedWallet

	currentNode := node_manager.CurrentNode
	if currentNode == nil {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return p.alertBox.Layout(gtx, th, lang.Translate("Unassigned node! Select your node from the node management page."))
		})
	} else {
		if walletapi.Connected && wallet != nil {
			nodeSynced := false
			walletSynced := false

			walletHeight := wallet.Memory.Get_Height()
			networkHeight := uint64(0)

			if currentNode.Integrated {
				nodeStatus := node_status_bar.Instance.IntegratedNodeStatus
				nodeHeight := uint64(nodeStatus.Height)
				networkHeight = uint64(nodeStatus.BestHeight)
				nodeSynced = nodeHeight >= networkHeight-8
				walletSynced = walletHeight >= networkHeight-8
			} else {
				nodeStatus := node_status_bar.Instance.RemoteNodeInfo.Result
				nodeHeight := uint64(nodeStatus.Height)
				networkHeight = uint64(nodeStatus.StableHeight)
				nodeSynced = nodeHeight >= networkHeight
				walletSynced = walletHeight >= networkHeight
			}

			if nodeSynced {
				isRegistered := wallet.Memory.IsRegistered()
				// check registration first because the wallet will never be synced if not registered
				if !isRegistered {
					widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
						return p.alertBox.Layout(gtx, th, lang.Translate("This wallet is not registered on the blockchain."))
					})

					widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top: unit.Dp(0), Bottom: unit.Dp(20),
							Left: unit.Dp(30), Right: unit.Dp(30),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							p.buttonRegister.Text = lang.Translate("REGISTER WALLET")
							return p.buttonRegister.Layout(gtx, th)
						})
					})
				} else if !walletSynced {
					widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
						text := lang.Translate("The wallet is not synced. Please wait and let it sync. The network height is currently {}.")
						return p.alertBox.Layout(gtx, th, strings.Replace(text, "{}", fmt.Sprint(networkHeight), -1))
					})
				}
			} else {
				widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
					text := lang.Translate("The node is out of synced. Please wait and let it sync. The network height is currently {}.")
					return p.alertBox.Layout(gtx, th, strings.Replace(text, "{}", fmt.Sprint(networkHeight), -1))
				})
			}
		} else {
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				return p.alertBox.Layout(gtx, th, lang.Translate("The wallet is not connected to a node."))
			})
		}
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Left: unit.Dp(30), Right: unit.Dp(30),
			Top: unit.Dp(0), Bottom: unit.Dp(30),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return p.displayBalance.Layout(gtx, th)
		})
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		rect := image.Rectangle{Max: image.Pt(gtx.Constraints.Max.X, 2)}
		paint.FillShape(gtx.Ops, color.NRGBA{A: 50}, clip.UniformRRect(rect, 0).Op(gtx.Ops))

		cl := clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Dp(1))}.Push(gtx.Ops)
		paint.ColorOp{Color: color.NRGBA{R: 0, G: 0, B: 0, A: 50}}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		cl.Pop()
		return layout.Dimensions{Size: rect.Max}
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(20), Bottom: unit.Dp(15),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			text := make(map[string]string)
			text["tokens"] = lang.Translate("Tokens")
			text["txs"] = lang.Translate("Transactions")
			return p.tabBars.Layout(gtx, th, text)
		})
	})

	{
		changed, tab := p.txBar.Changed()
		if changed {
			switch tab {
			case "all":
				p.getTransfersParams = wallet_manager.GetTransfersParams{}
			case "in":
				p.getTransfersParams = wallet_manager.GetTransfersParams{
					In: sql.NullBool{Bool: true, Valid: true},
				}
			case "out":
				p.getTransfersParams = wallet_manager.GetTransfersParams{
					Out: sql.NullBool{Bool: true, Valid: true},
				}
			case "coinbase":
				p.getTransfersParams = wallet_manager.GetTransfersParams{
					Coinbase: sql.NullBool{Bool: true, Valid: true},
				}
			}

			p.LoadTxs()
		}
	}

	{
		changed, key := p.tabBars.Changed()
		if changed {
			settings.App.TabBarsKey = key
			settings.Save()
		}
	}

	if p.tabBars.Key == "tokens" {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top: unit.Dp(0), Bottom: unit.Dp(15),
				Left: unit.Dp(30), Right: unit.Dp(30),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return p.tokenBar.Layout(gtx, th)
			})
		})

		if len(p.tokenItems) == 0 {
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top: unit.Dp(0), Bottom: unit.Dp(20),
					Left: unit.Dp(30), Right: unit.Dp(30),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), lang.Translate("You don't have any favorite tokens. Click the menu icon to manage tokens."))
					return lbl.Layout(gtx)
				})
			})
		}

		for i := range p.tokenItems {
			idx := i
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top: unit.Dp(0), Bottom: unit.Dp(15),
					Right: unit.Dp(30), Left: unit.Dp(30),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return p.tokenItems[idx].Layout(gtx, th)
				})
			})
		}
	}

	if p.tabBars.Key == "txs" {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top: unit.Dp(0), Bottom: unit.Dp(15),
				Left: unit.Dp(30), Right: unit.Dp(30),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return p.txBar.Layout(gtx, th)
			})
		})

		if len(p.txItems) == 0 {
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top: unit.Dp(0), Bottom: unit.Dp(20),
					Left: unit.Dp(30), Right: unit.Dp(30),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), lang.Translate("You don't have any txs. Try adjusting filering options or wait for wallet to sync."))
					return lbl.Layout(gtx)
				})
			})
		}

		for i := range p.txItems {
			idx := i
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top: unit.Dp(0), Bottom: unit.Dp(15),
					Right: unit.Dp(30), Left: unit.Dp(30),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return p.txItems[idx].Layout(gtx, th)
				})
			})
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
}

func NewAlertBox() *AlertBox {
	iconWarning, _ := widget.NewIcon(icons.AlertWarning)
	return &AlertBox{
		iconWarning: iconWarning,
	}
}

func (n *AlertBox) Layout(gtx layout.Context, th *material.Theme, text string) layout.Dimensions {
	border := widget.Border{
		Color:        color.NRGBA{A: 100},
		CornerRadius: unit.Dp(5),
		Width:        unit.Dp(1),
	}

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
						label := material.Label(th, unit.Sp(14), text)
						label.Color = color.NRGBA{A: 200}
						return label.Layout(gtx)
					}),
				)
			})
		})
	})
}

type SendReceiveButtons struct {
	ButtonSend    *components.Button
	ButtonReceive *components.Button
}

func NewSendReceiveButtons() *SendReceiveButtons {
	sendIcon, _ := widget.NewIcon(icons.NavigationArrowUpward)
	buttonSend := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
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

	return &SendReceiveButtons{
		ButtonSend:    buttonSend,
		ButtonReceive: buttonReceive,
	}
}

func (s *SendReceiveButtons) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.Y = gtx.Dp(40)
			s.ButtonSend.Text = lang.Translate("SEND")
			return s.ButtonSend.Layout(gtx, th)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(15)}.Layout),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.Y = gtx.Dp(40)
			s.ButtonReceive.Text = lang.Translate("RECEIVE")
			return s.ButtonReceive.Layout(gtx, th)
		}),
	)
}

type ButtonHideBalance struct {
	Button *components.Button

	hideBalanceIcon *widget.Icon
	showBalanceIcon *widget.Icon
}

func NewButtonHideBalance() *ButtonHideBalance {
	hideBalanceIcon, _ := widget.NewIcon(icons.ActionVisibility)
	showBalanceIcon, _ := widget.NewIcon(icons.ActionVisibilityOff)
	buttonHideBalance := components.NewButton(components.ButtonStyle{
		Icon:      hideBalanceIcon,
		TextColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
	})

	return &ButtonHideBalance{
		Button:          buttonHideBalance,
		hideBalanceIcon: hideBalanceIcon,
		showBalanceIcon: showBalanceIcon,
	}
}

func (b *ButtonHideBalance) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if settings.App.HideBalance {
		b.Button.Style.Icon = b.showBalanceIcon
	} else {
		b.Button.Style.Icon = b.hideBalanceIcon
	}

	if b.Button.Clicked() {
		settings.App.HideBalance = !settings.App.HideBalance
		settings.Save()
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	return b.Button.Layout(gtx, th)
}

type DisplayBalance struct {
	sendReceiveButtons *SendReceiveButtons
	buttonHideBalance  *ButtonHideBalance
}

func NewDisplayBalance() *DisplayBalance {
	sendReceiveButtons := NewSendReceiveButtons()
	buttonHideBalance := NewButtonHideBalance()

	return &DisplayBalance{
		buttonHideBalance:  buttonHideBalance,
		sendReceiveButtons: sendReceiveButtons,
	}
}

func (d *DisplayBalance) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	wallet := wallet_manager.OpenedWallet

	if d.sendReceiveButtons.ButtonSend.Clicked() {
		page_instance.pageSendForm.token = wallet_manager.DeroToken()
		page_instance.pageRouter.SetCurrent(PAGE_SEND_FORM)
		page_instance.header.AddHistory(PAGE_SEND_FORM)
	}

	if d.sendReceiveButtons.ButtonReceive.Clicked() {
		page_instance.pageRouter.SetCurrent(PAGE_RECEIVE_FORM)
		page_instance.header.AddHistory(PAGE_RECEIVE_FORM)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(14), lang.Translate("Available Balance"))
			lbl.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 150}

			return lbl.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					balance := uint64(0)
					if walletapi.Connected && wallet != nil {
						balance, _ = wallet.Memory.Get_Balance()
					}

					amount := utils.ShiftNumber{Number: balance, Decimals: 5}.Format()
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
					gtx.Constraints.Min.Y = gtx.Dp(30)
					gtx.Constraints.Min.X = gtx.Dp(30)
					return d.buttonHideBalance.Layout(gtx, th)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return d.sendReceiveButtons.Layout(gtx, th)
		}),
	)
}

type TokenBar struct {
	buttonListToken *components.Button
}

func NewTokenBar() *TokenBar {
	listIcon, _ := widget.NewIcon(icons.ActionViewList)
	buttonListToken := components.NewButton(components.ButtonStyle{
		Icon:            listIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{A: 255},
		Inset: layout.Inset{
			Top: unit.Dp(5), Bottom: unit.Dp(5),
			Left: unit.Dp(8), Right: unit.Dp(8),
		},
		Rounded:   components.UniformRounded(5),
		Animation: components.NewButtonAnimationDefault(),
	})

	return &TokenBar{
		buttonListToken: buttonListToken,
	}
}

func (t *TokenBar) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if t.buttonListToken.Clicked() {
		page_instance.pageRouter.SetCurrent(PAGE_SC_FOLDERS)
		page_instance.header.AddHistory(PAGE_SC_FOLDERS)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(18), lang.Translate("Favorites"))
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return t.buttonListToken.Layout(gtx, th)
				}),
			)
		}),
	)
}

type TokenListItem struct {
	token     *wallet_manager.Token
	clickable *widget.Clickable
	image     *prefabs.ImageHoverClick

	balance uint64
}

func NewTokenListItem(token wallet_manager.Token, img paint.ImageOp) *TokenListItem {
	return &TokenListItem{
		token:     &token,
		image:     prefabs.NewImageHoverClick(img),
		clickable: new(widget.Clickable),
	}
}

func (item *TokenListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if item.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
	}

	if item.clickable.Clicked() {
		page_instance.pageSCToken.token = item.token
		page_instance.pageRouter.SetCurrent(PAGE_SC_TOKEN)
		page_instance.header.AddHistory(PAGE_SC_TOKEN)
	}

	m := op.Record(gtx.Ops)
	dims := item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(13), Bottom: unit.Dp(13),
			Left: unit.Dp(15), Right: unit.Dp(15),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = gtx.Dp(50)
					gtx.Constraints.Max.Y = gtx.Dp(50)
					return item.image.Layout(gtx)
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
									lbl := material.Label(th, unit.Sp(18), item.token.Name)
									lbl.Font.Weight = font.Bold
									return lbl.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									scId := utils.ReduceTxId(item.token.SCID)
									if item.token.Symbol.Valid {
										scId = fmt.Sprintf("%s (%s)", scId, item.token.Symbol.String)
									}

									lbl := material.Label(th, unit.Sp(14), scId)
									lbl.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 150}
									return lbl.Layout(gtx)
								}),
							)
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
			NW:   gtx.Dp(10), NE: gtx.Dp(10),
			SE: gtx.Dp(10), SW: gtx.Dp(10),
		}.Op(gtx.Ops))

	c.Add(gtx.Ops)

	if !settings.App.HideBalance {
		layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			r := op.Record(gtx.Ops)
			labelDims := layout.Inset{
				Left: unit.Dp(8), Right: unit.Dp(8),
				Bottom: unit.Dp(5), Top: unit.Dp(5),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				balance := utils.ShiftNumber{Number: uint64(item.balance), Decimals: int(item.token.Decimals)}
				label := material.Label(th, unit.Sp(18), balance.Format())
				label.Font.Weight = font.Bold
				return label.Layout(gtx)
			})
			c := r.Stop()

			x := float32(gtx.Dp(5))
			y := float32(dims.Size.Y/2 - labelDims.Size.Y/2)
			offset := f32.Affine2D{}.Offset(f32.Pt(x, y))
			defer op.Affine(offset).Push(gtx.Ops).Pop()

			paint.FillShape(gtx.Ops, color.NRGBA{R: 225, G: 225, B: 225, A: 255},
				clip.RRect{
					Rect: image.Rectangle{Max: labelDims.Size},
					NW:   gtx.Dp(5), NE: gtx.Dp(5),
					SE: gtx.Dp(5), SW: gtx.Dp(5),
				}.Op(gtx.Ops))

			c.Add(gtx.Ops)
			return labelDims
		})
	}

	return dims
}

type TxBar struct {
	buttonAll      *components.Button
	buttonIn       *components.Button
	buttonOut      *components.Button
	buttonCoinbase *components.Button
	buttonFilter   *components.Button
	txCount        int

	textColorOn  color.NRGBA
	textColorOff color.NRGBA
	bgColorOn    color.NRGBA
	bgColorOff   color.NRGBA

	tab     string
	changed bool
}

func NewTxBar() *TxBar {
	buttonAll := components.NewButton(components.ButtonStyle{
		TextSize: unit.Sp(16),
		Inset: layout.Inset{
			Top: unit.Dp(5), Bottom: unit.Dp(5),
			Left: unit.Dp(8), Right: unit.Dp(8),
		},
		Rounded:   components.UniformRounded(5),
		Animation: components.NewButtonAnimationDefault(),
	})

	buttonIn := components.NewButton(components.ButtonStyle{
		TextSize: unit.Sp(16),
		Inset: layout.Inset{
			Top: unit.Dp(5), Bottom: unit.Dp(5),
			Left: unit.Dp(8), Right: unit.Dp(8),
		},
		Rounded:   components.UniformRounded(5),
		Animation: components.NewButtonAnimationDefault(),
	})

	buttonOut := components.NewButton(components.ButtonStyle{
		TextSize: unit.Sp(16),
		Inset: layout.Inset{
			Top: unit.Dp(5), Bottom: unit.Dp(5),
			Left: unit.Dp(8), Right: unit.Dp(8),
		},
		Rounded:   components.UniformRounded(5),
		Animation: components.NewButtonAnimationDefault(),
	})

	buttonCoinbase := components.NewButton(components.ButtonStyle{
		TextSize: unit.Sp(16),
		Inset: layout.Inset{
			Top: unit.Dp(5), Bottom: unit.Dp(5),
			Left: unit.Dp(8), Right: unit.Dp(8),
		},
		Rounded:   components.UniformRounded(5),
		Animation: components.NewButtonAnimationDefault(),
	})

	buttonFilter := components.NewButton(components.ButtonStyle{
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{A: 255},
		TextSize:        unit.Sp(16),
		Inset: layout.Inset{
			Top: unit.Dp(5), Bottom: unit.Dp(5),
			Left: unit.Dp(8), Right: unit.Dp(8),
		},
		Rounded:   components.UniformRounded(5),
		Animation: components.NewButtonAnimationDefault(),
	})

	textColorOn := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	textColorOff := color.NRGBA{A: 255}

	bgColorOn := color.NRGBA{A: 255}
	bgColorOff := color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	return &TxBar{
		buttonAll:      buttonAll,
		buttonIn:       buttonIn,
		buttonOut:      buttonOut,
		buttonCoinbase: buttonCoinbase,
		buttonFilter:   buttonFilter,
		tab:            "all",

		textColorOn:  textColorOn,
		textColorOff: textColorOff,
		bgColorOn:    bgColorOn,
		bgColorOff:   bgColorOff,
	}
}

func (t *TxBar) Changed() (bool, string) {
	return t.changed, t.tab
}

func (t *TxBar) setActiveButton(button *components.Button, tab string) {
	if t.tab == tab {
		button.Style.TextColor = t.textColorOn
		button.Style.BackgroundColor = t.bgColorOn
		button.Disabled = true
	} else {
		button.Style.TextColor = t.textColorOff
		button.Style.BackgroundColor = t.bgColorOff
		button.Disabled = false
	}
}

func (t *TxBar) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	t.changed = false

	if t.buttonAll.Clicked() {
		t.changed = true
		t.tab = "all"
	}

	if t.buttonIn.Clicked() {
		t.changed = true
		t.tab = "in"
	}

	if t.buttonOut.Clicked() {
		t.changed = true
		t.tab = "out"
	}

	if t.buttonCoinbase.Clicked() {
		t.changed = true
		t.tab = "coinbase"
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							t.buttonAll.Text = lang.Translate("All")
							t.setActiveButton(t.buttonAll, "all")
							return t.buttonAll.Layout(gtx, th)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							t.buttonIn.Text = lang.Translate("In")
							t.setActiveButton(t.buttonIn, "in")
							return t.buttonIn.Layout(gtx, th)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							t.buttonOut.Text = lang.Translate("Out")
							t.setActiveButton(t.buttonOut, "out")
							return t.buttonOut.Layout(gtx, th)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							t.buttonCoinbase.Text = lang.Translate("Coinbase")
							t.setActiveButton(t.buttonCoinbase, "coinbase")
							return t.buttonCoinbase.Layout(gtx, th)
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					t.buttonFilter.Text = lang.Translate("Filter")
					return t.buttonFilter.Layout(gtx, th)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			txt := lang.Translate("{} transactions")
			txt = strings.Replace(txt, "{}", fmt.Sprint(t.txCount), -1)
			lbl := material.Label(th, unit.Sp(14), txt)
			return lbl.Layout(gtx)
		}),
	)
}

type TxListItem struct {
	entry     rpc.Entry
	clickable *widget.Clickable
	image     *components.Image
	decimals  int
}

func NewTxListItem(entry rpc.Entry, img paint.ImageOp) *TxListItem {
	return &TxListItem{
		entry: entry,
		image: &components.Image{
			Src: img,
			Fit: components.Cover,
		},
		clickable: new(widget.Clickable),
		decimals:  5,
	}
}

func (item *TxListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if item.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
	}

	if item.clickable.Clicked() {
		page_instance.pageTransaction.entry = &item.entry
		page_instance.pageRouter.SetCurrent(PAGE_TRANSACTION)
		page_instance.header.AddHistory(PAGE_TRANSACTION)
	}

	m := op.Record(gtx.Ops)
	dims := item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(13), Bottom: unit.Dp(13),
			Left: unit.Dp(15), Right: unit.Dp(15),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = gtx.Dp(35)
					gtx.Constraints.Max.Y = gtx.Dp(35)
					return item.image.Layout(gtx)
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
									txt := ""
									if item.entry.Coinbase {
										txt = lang.Translate("From Coinbase")
									} else {
										txt = utils.ReduceTxId(item.entry.TXID)
									}

									lbl := material.Label(th, unit.Sp(18), txt)
									lbl.Font.Weight = font.Bold
									return lbl.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									txt := item.entry.Time.Format("2006-01-02 15:04")
									lbl := material.Label(th, unit.Sp(14), txt)
									lbl.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 150}
									return lbl.Layout(gtx)
								}),
							)
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
			NW:   gtx.Dp(10), NE: gtx.Dp(10),
			SE: gtx.Dp(10), SW: gtx.Dp(10),
		}.Op(gtx.Ops))

	c.Add(gtx.Ops)

	if !settings.App.HideBalance {
		layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			r := op.Record(gtx.Ops)
			labelDims := layout.Inset{
				Left: unit.Dp(8), Right: unit.Dp(8),
				Bottom: unit.Dp(5), Top: unit.Dp(5),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				amount := item.entry.Amount
				balance := utils.ShiftNumber{Number: amount, Decimals: item.decimals}
				label := material.Label(th, unit.Sp(18), balance.Format())
				label.Font.Weight = font.Bold
				return label.Layout(gtx)
			})
			c := r.Stop()

			x := float32(gtx.Dp(5))
			y := float32(dims.Size.Y/2 - labelDims.Size.Y/2)
			offset := f32.Affine2D{}.Offset(f32.Pt(x, y))
			defer op.Affine(offset).Push(gtx.Ops).Pop()

			paint.FillShape(gtx.Ops, color.NRGBA{R: 225, G: 225, B: 225, A: 255},
				clip.RRect{
					Rect: image.Rectangle{Max: labelDims.Size},
					NW:   gtx.Dp(5), NE: gtx.Dp(5),
					SE: gtx.Dp(5), SW: gtx.Dp(5),
				}.Op(gtx.Ops))

			c.Add(gtx.Ops)
			return labelDims
		})
	}

	return dims
}
