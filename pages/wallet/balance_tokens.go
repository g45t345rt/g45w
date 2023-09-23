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
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_icons"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/image_modal"
	"github.com/g45t345rt/g45w/containers/node_status_bar"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/node_manager"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/theme"
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

	alertBox         *AlertBox
	displayBalance   *DisplayBalance
	tokenBar         *TokenBar
	tokenItems       []*TokenListItem
	buttonSettings   *components.Button
	buttonRegister   *components.Button
	buttonCopyAddr   *components.Button
	buttonDexSwap    *components.Button
	tabBars          *components.TabBars
	txBar            *TxBar
	txItems          []*TxListItem
	getEntriesParams wallet_manager.GetEntriesParams
	tokenDragItems   *components.DragItems
	tokenList        *widget.List

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
		Animation: components.NewButtonAnimationScale(.98),
	})

	registerIcon, _ := widget.NewIcon(icons.ActionAssignmentTurnedIn)
	buttonRegister := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		TextSize:  unit.Sp(14),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
		Icon:      registerIcon,
		IconGap:   unit.Dp(10),
	})
	buttonRegister.Label.Alignment = text.Middle
	buttonRegister.Style.Font.Weight = font.Bold

	swapIcon, _ := widget.NewIcon(app_icons.Swap)
	buttonDexSwap := components.NewButton(components.ButtonStyle{
		Icon:      swapIcon,
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
	buttonDexSwap.Label.Alignment = text.Middle
	buttonDexSwap.Style.Font.Weight = font.Bold

	copyIcon, _ := widget.NewIcon(icons.ContentContentCopy)
	buttonCopyAddr := components.NewButton(components.ButtonStyle{
		Icon: copyIcon,
	})

	tabBarsItems := []*components.TabBarsItem{
		components.NewTabBarItem("txs"),
		components.NewTabBarItem("tokens"),
	}
	defaultTabKey := settings.App.MainTabBars
	tabBars := components.NewTabBars(defaultTabKey, tabBarsItems)

	txBar := NewTxBar()
	tokenDragItems := components.NewDragItems()
	tokenList := new(widget.List)
	tokenList.Axis = layout.Vertical

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
		buttonDexSwap:  buttonDexSwap,
		tabBars:        tabBars,
		txBar:          txBar,
		tokenDragItems: tokenDragItems,
		tokenList:      tokenList,
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
	for _, token := range tokens {
		tokenItems = append(tokenItems, NewTokenListItem(token))
	}

	p.tokenItems = tokenItems
	return nil
}

func (p *PageBalanceTokens) LoadTxs() {
	wallet := wallet_manager.OpenedWallet
	entries := wallet.GetEntries(&crypto.ZEROHASH, p.getEntriesParams)

	txItems := []*TxListItem{}

	for _, entry := range entries {
		txItems = append(txItems, NewTxListItem(entry, 5))
	}

	p.txItems = txItems
	p.txBar.txCount = len(entries)
}

func (p *PageBalanceTokens) ResetWalletHeader() {
	wallet := wallet_manager.OpenedWallet
	page_instance.header.Title = func() string {
		return fmt.Sprintf("%s [%s]", lang.Translate("Wallet"), wallet.Info.Name)
	}

	page_instance.header.ButtonRight = nil
	addr := wallet.Memory.GetAddress().String()
	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		if p.buttonCopyAddr.Clicked() {
			clipboard.WriteOp{
				Text: addr,
			}.Add(gtx.Ops)
			notification_modals.InfoInstance.SetText(lang.Translate("Clipboard"), lang.Translate("Addr copied to clipboard"))
			notification_modals.InfoInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		}

		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				walletAddr := utils.ReduceAddr(addr)
				label := material.Label(th, unit.Sp(16), walletAddr)
				label.Color = theme.Current.TextMuteColor
				return label.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(5)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X = gtx.Dp(18)
				gtx.Constraints.Max.Y = gtx.Dp(18)
				p.buttonCopyAddr.Style.Colors = theme.Current.ModalButtonColors
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

	if p.buttonDexSwap.Clicked() {
		page_instance.pageRouter.SetCurrent(PAGE_DEX_PAIRS)
		page_instance.header.AddHistory(PAGE_DEX_PAIRS)
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
			stableHeight := uint64(0)
			nodeHeight := uint64(0)

			if currentNode.Integrated {
				nodeStatus := node_status_bar.Instance.IntegratedNodeStatus
				nodeHeight = uint64(nodeStatus.Height)
				stableHeight = uint64(nodeStatus.BestHeight)
				nodeSynced = nodeHeight >= stableHeight-8
				walletSynced = walletHeight >= stableHeight-8
			} else {
				nodeStatus := node_status_bar.Instance.RemoteNodeInfo.Result
				nodeHeight = uint64(nodeStatus.Height)
				stableHeight = uint64(nodeStatus.StableHeight)
				nodeSynced = nodeHeight >= stableHeight
				walletSynced = walletHeight >= stableHeight
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
							p.buttonRegister.Style.Colors = theme.Current.ButtonPrimaryColors
							return p.buttonRegister.Layout(gtx, th)
						})
					})
				} else if !walletSynced {
					widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
						text := lang.Translate("The wallet is not synced. Please wait and let it sync. The node height is currently {}.")
						return p.alertBox.Layout(gtx, th, strings.Replace(text, "{}", fmt.Sprint(nodeHeight), -1))
					})
				}
			} else {
				widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
					text := lang.Translate("The node is out of synced. Please wait and let it sync. The stable height is currently {}.")
					return p.alertBox.Layout(gtx, th, strings.Replace(text, "{}", fmt.Sprint(stableHeight), -1))
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
			Top: unit.Dp(0), Bottom: unit.Dp(20),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return p.displayBalance.Layout(gtx, th)
		})
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Left: unit.Dp(30), Right: unit.Dp(30),
			Top: unit.Dp(0), Bottom: unit.Dp(30),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			p.buttonDexSwap.Style.Colors = theme.Current.ButtonSecondaryColors
			p.buttonDexSwap.Text = lang.Translate("DEX Swap")
			return p.buttonDexSwap.Layout(gtx, th)
		})
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return prefabs.Divider(gtx, 3)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(20), Bottom: unit.Dp(15),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			text := make(map[string]string)
			text["tokens"] = lang.Translate("Tokens")
			text["txs"] = lang.Translate("Transactions")
			p.tabBars.Colors = theme.Current.TabBarsColors
			return p.tabBars.Layout(gtx, th, unit.Sp(18), text)
		})
	})

	{
		changed, tab := p.txBar.Changed()
		if changed {
			switch tab {
			case "all":
				p.getEntriesParams = wallet_manager.GetEntriesParams{}
			case "in":
				p.getEntriesParams = wallet_manager.GetEntriesParams{
					In: sql.NullBool{Bool: true, Valid: true},
				}
			case "out":
				p.getEntriesParams = wallet_manager.GetEntriesParams{
					Out: sql.NullBool{Bool: true, Valid: true},
				}
			case "coinbase":
				p.getEntriesParams = wallet_manager.GetEntriesParams{
					Coinbase: sql.NullBool{Bool: true, Valid: true},
				}
			}

			p.LoadTxs()
		}
	}

	{
		changed, key := p.tabBars.Changed()
		if changed {
			settings.App.MainTabBars = key
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
					lbl := material.Label(th, unit.Sp(16), lang.Translate("You don't have any favorite tokens. Click the folder icon to manage tokens."))
					lbl.Color = theme.Current.TextMuteColor
					return lbl.Layout(gtx)
				})
			})
		}

		/*
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
			}*/
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return p.tokenDragItems.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
				return p.tokenList.List.Layout(gtx, len(p.tokenItems), func(gtx layout.Context, index int) layout.Dimensions {
					p.tokenDragItems.LayoutItem(gtx, index, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top: unit.Dp(0), Bottom: unit.Dp(15),
							Right: unit.Dp(30), Left: unit.Dp(30),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return p.tokenItems[index].Layout(gtx, th)
						})
					})

					return layout.Inset{
						Top: unit.Dp(0), Bottom: unit.Dp(15),
						Right: unit.Dp(30), Left: unit.Dp(30),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return p.tokenItems[index].Layout(gtx, th)
					})
				})
			})
		})
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
					lbl := material.Label(th, unit.Sp(16), lang.Translate("You don't have any txs. Try adjusting filtering options or wait for wallet to sync."))
					lbl.Color = theme.Current.TextMuteColor
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
		return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
	})

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return widgets[index](gtx)
	})
}

type AlertBox struct {
	iconWarning *widget.Icon
	textEditor  *widget.Editor
}

func NewAlertBox() *AlertBox {
	iconWarning, _ := widget.NewIcon(icons.AlertWarning)
	textEditor := new(widget.Editor)
	textEditor.ReadOnly = true

	return &AlertBox{
		iconWarning: iconWarning,
		textEditor:  textEditor,
	}
}

func (n *AlertBox) Layout(gtx layout.Context, th *material.Theme, text string) layout.Dimensions {
	color := theme.Current.TextMuteColor
	border := widget.Border{
		Color:        color,
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
						return n.iconWarning.Layout(gtx, color)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						editor := material.Editor(th, n.textEditor, "")
						editor.Color = color
						editor.TextSize = unit.Sp(14)

						if n.textEditor.Text() != text {
							n.textEditor.SetText(text)
						}

						return editor.Layout(gtx)
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
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      sendIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonSend.Label.Alignment = text.Middle
	buttonSend.Style.Font.Weight = font.Bold

	receiveIcon, _ := widget.NewIcon(icons.NavigationArrowDownward)
	buttonReceive := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      receiveIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
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
			s.ButtonSend.Style.Colors = theme.Current.ButtonPrimaryColors
			return s.ButtonSend.Layout(gtx, th)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(15)}.Layout),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.Y = gtx.Dp(40)
			s.ButtonReceive.Text = lang.Translate("RECEIVE")
			s.ButtonReceive.Style.Colors = theme.Current.ButtonPrimaryColors
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
		Icon: hideBalanceIcon,
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
	balanceEditor      *widget.Editor
}

func NewDisplayBalance() *DisplayBalance {
	sendReceiveButtons := NewSendReceiveButtons()
	buttonHideBalance := NewButtonHideBalance()

	balanceEditor := new(widget.Editor)
	balanceEditor.ReadOnly = true
	balanceEditor.SingleLine = true

	return &DisplayBalance{
		buttonHideBalance:  buttonHideBalance,
		sendReceiveButtons: sendReceiveButtons,
		balanceEditor:      balanceEditor,
	}
}

func (d *DisplayBalance) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	wallet := wallet_manager.OpenedWallet

	if d.sendReceiveButtons.ButtonSend.Clicked() {
		page_instance.pageSendForm.SetToken(wallet_manager.DeroToken())
		page_instance.pageSendForm.ClearForm()
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
			lbl.Color = theme.Current.TextMuteColor

			return lbl.Layout(gtx)
		}),
		//layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					balance, _ := wallet.Memory.Get_Balance()
					amount := utils.ShiftNumber{Number: balance, Decimals: 5}.Format()

					if d.balanceEditor.Text() != amount {
						d.balanceEditor.SetText(amount)
					}

					r := op.Record(gtx.Ops)
					amountEditor := material.Editor(th, d.balanceEditor, "")
					amountEditor.TextSize = unit.Sp(34)
					amountEditor.Font.Weight = font.Bold

					dims := amountEditor.Layout(gtx)
					c := r.Stop()

					if settings.App.HideBalance {
						paint.FillShape(gtx.Ops, theme.Current.HideBalanceBgColor, clip.Rect{
							Max: dims.Size,
						}.Op())
					} else {
						c.Add(gtx.Ops)
					}

					return dims
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.Y = gtx.Dp(30)
					gtx.Constraints.Min.X = gtx.Dp(30)
					d.buttonHideBalance.Button.Style.Colors = theme.Current.ButtonIconPrimaryColors
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
	buttonManageTokens *components.Button
}

func NewTokenBar() *TokenBar {
	folderIcon, _ := widget.NewIcon(icons.FileFolder)

	buttonManageTokens := components.NewButton(components.ButtonStyle{
		Icon: folderIcon,
		Inset: layout.Inset{
			Top: unit.Dp(5), Bottom: unit.Dp(5),
			Left: unit.Dp(8), Right: unit.Dp(8),
		},
		Rounded:   components.UniformRounded(5),
		Animation: components.NewButtonAnimationDefault(),
	})

	return &TokenBar{
		buttonManageTokens: buttonManageTokens,
	}
}

func (t *TokenBar) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if t.buttonManageTokens.Clicked() {
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
					t.buttonManageTokens.Style.Colors = theme.Current.ButtonPrimaryColors
					return t.buttonManageTokens.Layout(gtx, th)
				}),
			)
		}),
	)
}

type TokenListItem struct {
	token      *wallet_manager.Token
	clickable  *widget.Clickable
	imageHover *prefabs.ImageHoverClick
}

func NewTokenListItem(token wallet_manager.Token) *TokenListItem {
	return &TokenListItem{
		token:      &token,
		imageHover: prefabs.NewImageHoverClick(),
		clickable:  new(widget.Clickable),
	}
}

func (item *TokenListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if item.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
	}

	if item.clickable.Clicked() {
		page_instance.pageSCToken.SetToken(item.token)
		page_instance.pageRouter.SetCurrent(PAGE_SC_TOKEN)
		page_instance.header.AddHistory(PAGE_SC_TOKEN)
	}

	m := op.Record(gtx.Ops)

	dims := layout.Inset{
		Top: unit.Dp(13), Bottom: unit.Dp(13),
		Left: unit.Dp(15), Right: unit.Dp(15),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if item.imageHover.Clickable.Clicked() {
					image_modal.Instance.Open(item.token.Name, item.imageHover.Image.Src)
				}

				item.imageHover.Image.Src = item.token.LoadImageOp()
				gtx.Constraints.Max.X = gtx.Dp(50)
				gtx.Constraints.Max.Y = gtx.Dp(50)
				return item.imageHover.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									lbl := material.Label(th, unit.Sp(18), item.token.Name)
									if len(item.token.Name) > 20 {
										lbl.TextSize = unit.Sp(16)
									}

									lbl.Font.Weight = font.Bold
									return lbl.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									scId := utils.ReduceTxId(item.token.SCID)
									if item.token.Symbol.Valid {
										scId = fmt.Sprintf("%s (%s)", scId, item.token.Symbol.String)
									}

									lbl := material.Label(th, unit.Sp(14), scId)
									lbl.Color = theme.Current.TextMuteColor
									return lbl.Layout(gtx)
								}),
							)
						}),
					)
				})
			}),
		)
	})
	c := m.Stop()

	paint.FillShape(gtx.Ops, theme.Current.ListBgColor,
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
				wallet := wallet_manager.OpenedWallet
				balance, _ := wallet.Memory.Get_Balance_scid(item.token.GetHash())
				amount := utils.ShiftNumber{Number: uint64(balance), Decimals: int(item.token.Decimals)}
				lbl := material.Label(th, unit.Sp(18), amount.Format())
				lbl.Font.Weight = font.Bold
				return lbl.Layout(gtx)
			})
			c := r.Stop()

			x := float32(gtx.Dp(5))
			y := float32(dims.Size.Y/2 - labelDims.Size.Y/2)
			offset := f32.Affine2D{}.Offset(f32.Pt(x, y))
			defer op.Affine(offset).Push(gtx.Ops).Pop()

			paint.FillShape(gtx.Ops, theme.Current.ListItemTagBgColor,
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

	filterIcon, _ := widget.NewIcon(app_icons.Filter)
	buttonFilter := components.NewButton(components.ButtonStyle{
		//TextSize: unit.Sp(16),
		Icon: filterIcon,
		Inset: layout.Inset{
			Top: unit.Dp(8), Bottom: unit.Dp(8),
			Left: unit.Dp(8), Right: unit.Dp(8),
		},
		Border: widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Width:        unit.Dp(1),
			CornerRadius: unit.Dp(5),
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
		button.Style.Colors = theme.Current.ButtonPrimaryColors
		button.Disabled = true
	} else {
		button.Style.Colors = theme.Current.ButtonInvertColors
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
					//t.buttonFilter.Text = lang.Translate("Filter")
					gtx.Constraints.Max = image.Pt(gtx.Dp(30), gtx.Dp(30))
					t.buttonFilter.Style.Colors = theme.Current.ButtonSecondaryColors
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
	entry     wallet_manager.Entry
	clickable *widget.Clickable
	image     *components.Image
	decimals  int
}

func NewTxListItem(entry wallet_manager.Entry, decimals int) *TxListItem {
	return &TxListItem{
		entry: entry,
		image: &components.Image{
			Fit: components.Cover,
		},
		clickable: new(widget.Clickable),
		decimals:  decimals,
	}
}

func (item *TxListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if item.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
	}

	if item.clickable.Clicked() {
		page_instance.pageTransaction.SetEntry(item.entry)
		page_instance.pageRouter.SetCurrent(PAGE_TRANSACTION)
		page_instance.header.AddHistory(PAGE_TRANSACTION)
	}

	if item.entry.Incoming {
		item.image.Src = theme.Current.ArrowDownArcImage
	} else {
		item.image.Src = theme.Current.ArrowUpArcImage
	}

	if item.entry.Coinbase {
		item.image.Src = theme.Current.CoinbaseImage
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
									lbl.Color = theme.Current.TextMuteColor
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

	paint.FillShape(gtx.Ops, theme.Current.ListBgColor,
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

			paint.FillShape(gtx.Ops, theme.Current.ListItemTagBgColor,
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
