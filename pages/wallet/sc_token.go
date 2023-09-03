package page_wallet

import (
	"context"
	"database/sql"
	"fmt"
	"image"
	"strconv"

	"gioui.org/font"
	"gioui.org/io/clipboard"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/build_tx_modal"
	"github.com/g45t345rt/g45w/containers/confirm_modal"
	"github.com/g45t345rt/g45w/containers/image_modal"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/containers/prompt_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/sc"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageSCToken struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonOpenMenu      *components.Button
	tokenMenuSelect     *TokenMenuSelect
	sendReceiveButtons  *SendReceiveButtons
	tabBars             *components.TabBars
	txBar               *TxBar
	getTransfersParams  wallet_manager.GetTransfersParams
	txItems             []*TxListItem
	tokenInfo           *TokenInfoList
	balanceContainer    *BalanceContainer
	g45DisplayContainer *G45DisplayContainer
	buttonCopySCID      *components.Button

	token      *wallet_manager.Token
	tokenImage *prefabs.ImageHoverClick
	scIdEditor *widget.Editor

	list *widget.List
}

var _ router.Page = &PageSCToken{}

func NewPageSCToken() *PageSCToken {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	addIcon, _ := widget.NewIcon(icons.NavigationMenu)
	buttonOpenMenu := components.NewButton(components.ButtonStyle{
		Icon:      addIcon,
		Animation: components.NewButtonAnimationScale(.98),
	})

	list := new(widget.List)
	list.Axis = layout.Vertical

	image := prefabs.NewImageHoverClick()

	scIdEditor := new(widget.Editor)
	scIdEditor.WrapPolicy = text.WrapGraphemes
	scIdEditor.ReadOnly = true

	sendReceiveButtons := NewSendReceiveButtons()

	tabBarsItems := []*components.TabBarsItem{
		components.NewTabBarItem("txs"),
		components.NewTabBarItem("info"),
	}

	tabBars := components.NewTabBars("txs", tabBarsItems)
	txBar := NewTxBar()

	balanceContainer := NewBalanceContainer()
	g45DisplayContainer := NewG45DisplayContainer()

	copyIcon, _ := widget.NewIcon(icons.ContentContentCopy)
	buttonCopySCID := components.NewButton(components.ButtonStyle{
		Icon: copyIcon,
	})

	return &PageSCToken{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonOpenMenu:      buttonOpenMenu,
		tokenMenuSelect:     NewTokenMenuSelect(),
		tokenImage:          image,
		scIdEditor:          scIdEditor,
		sendReceiveButtons:  sendReceiveButtons,
		tabBars:             tabBars,
		txBar:               txBar,
		balanceContainer:    balanceContainer,
		g45DisplayContainer: g45DisplayContainer,
		buttonCopySCID:      buttonCopySCID,

		list: list,
	}
}

func (p *PageSCToken) IsActive() bool {
	return p.isActive
}

func (p *PageSCToken) Enter() {
	p.isActive = true

	wallet := wallet_manager.OpenedWallet
	wallet.Memory.TokenAdd(p.token.GetHash()) // we don't check error because the only possible error is if the token was already added

	p.tokenInfo = NewTokenInfoList(p.token)

	page_instance.header.Title = func() string { return p.token.Name }
	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		if p.buttonCopySCID.Clicked() {
			clipboard.WriteOp{
				Text: p.token.SCID,
			}.Add(gtx.Ops)
			notification_modals.InfoInstance.SetText(lang.Translate("Clipboard"), lang.Translate("Smart Contract ID copied to clipboard"))
			notification_modals.InfoInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		}

		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				scId := utils.ReduceTxId(p.token.SCID)
				if p.token.Symbol.Valid {
					scId = fmt.Sprintf("%s (%s)", scId, p.token.Symbol.String)
				}

				lbl := material.Label(th, unit.Sp(16), scId)
				lbl.Color = theme.Current.TextMuteColor
				return lbl.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(5)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X = gtx.Dp(18)
				gtx.Constraints.Max.Y = gtx.Dp(18)
				p.buttonCopySCID.Style.Colors = theme.Current.ModalButtonColors
				return p.buttonCopySCID.Layout(gtx, th)
			}),
		)
	}
	page_instance.header.ButtonRight = p.buttonOpenMenu
	p.scIdEditor.SetText(p.token.SCID)

	if !page_instance.header.IsHistory(PAGE_SC_TOKEN) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageSCToken) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageSCToken) LoadTxs() {
	wallet := wallet_manager.OpenedWallet
	entries := wallet.GetTransfers(p.token.SCID, p.getTransfersParams)

	txItems := []*TxListItem{}

	for _, entry := range entries {
		txItems = append(txItems, NewTxListItem(entry, int(p.token.Decimals)))
	}

	p.txItems = txItems
	p.txBar.txCount = len(entries)
}

func (p *PageSCToken) SetToken(token *wallet_manager.Token) {
	p.token = token
	p.token.RefreshImageOp()
	p.balanceContainer.SetToken(p.token)
	p.g45DisplayContainer.SetToken(p.token)
	p.g45DisplayContainer.Load()
	p.LoadTxs()
}

func (p *PageSCToken) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonOpenMenu.Clickable.Clicked() {
		p.tokenMenuSelect.SelectModal.Modal.SetVisible(true)
	}

	if p.tokenImage.Clickable.Clicked() {
		image_modal.Instance.Open(p.token.Name, p.tokenImage.Image.Src)
	}

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

	selected, sKey := p.tokenMenuSelect.SelectModal.Selected()
	if selected {
		wallet := wallet_manager.OpenedWallet
		var err error
		var successMsg = ""

		switch sKey {
		case "refresh_cache":
			wallet.ResetBalanceResult(p.token.SCID)
			successMsg = lang.Translate("Cache refreshed.")
		case "add_favorite":
			p.token.IsFavorite = sql.NullBool{Bool: true, Valid: true}
			err = wallet.UpdateToken(*p.token)
			successMsg = lang.Translate("Token added to favorites.")
		case "remove_favorite":
			p.token.IsFavorite = sql.NullBool{Bool: false, Valid: true}
			err = wallet.UpdateToken(*p.token)
			successMsg = lang.Translate("Token removed from favorites.")
		case "remove_token":
			go func() {
				yesChan := confirm_modal.Instance.Open(confirm_modal.ConfirmText{})

				for yes := range yesChan {
					if yes {
						wallet := wallet_manager.OpenedWallet
						err := wallet.DelToken(p.token.ID)

						if err != nil {
							notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
							notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
						} else {
							page_instance.header.GoBack()
							notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("Token removed."))
							notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
							p.tokenMenuSelect.SelectModal.Modal.SetVisible(false)
						}
					}
				}
			}()
		case "g45_display_nft":
			go func() {
				scId := p.token.GetHash()
				build_tx_modal.Instance.OpenWithRandomAddr(scId, func(randomAddr string, open func(txPayload build_tx_modal.TxPayload)) {
					open(build_tx_modal.TxPayload{
						Transfers: []rpc.Transfer{
							{SCID: scId, Destination: randomAddr, Burn: uint64(1)},
						},
						Ringsize: 2,
						SCArgs: rpc.Arguments{
							{Name: rpc.SCACTION, DataType: rpc.DataUint64, Value: uint64(rpc.SC_CALL)},
							{Name: rpc.SCID, DataType: rpc.DataHash, Value: scId},
							{Name: "entrypoint", DataType: rpc.DataString, Value: "DisplayNFT"},
						},
					})
				})
			}()
		case "g45_retrieve_nft":
			go func() {
				scId := p.token.GetHash()
				build_tx_modal.Instance.Open(build_tx_modal.TxPayload{
					Ringsize: 2,
					SCArgs: rpc.Arguments{
						{Name: rpc.SCACTION, DataType: rpc.DataUint64, Value: uint64(rpc.SC_CALL)},
						{Name: rpc.SCID, DataType: rpc.DataHash, Value: scId},
						{Name: "entrypoint", DataType: rpc.DataString, Value: "RetrieveNFT"},
					},
				})
			}()
		case "g45_display_token":
			go func() {
				txtChan := prompt_modal.Instance.Open("", lang.Translate("Enter amount"), key.HintNumeric)
				for txt := range txtChan {
					amount := utils.ShiftNumber{Decimals: int(p.token.Decimals)}
					err := amount.Parse(txt)
					if err != nil {
						return
					}

					scId := p.token.GetHash()
					build_tx_modal.Instance.OpenWithRandomAddr(scId, func(randomAddr string, open func(txPayload build_tx_modal.TxPayload)) {
						open(build_tx_modal.TxPayload{
							Transfers: []rpc.Transfer{
								{SCID: scId, Destination: randomAddr, Burn: amount.Number},
							},
							Ringsize: 2,
							SCArgs: rpc.Arguments{
								{Name: rpc.SCACTION, DataType: rpc.DataUint64, Value: uint64(rpc.SC_CALL)},
								{Name: rpc.SCID, DataType: rpc.DataHash, Value: scId},
								{Name: "entrypoint", DataType: rpc.DataString, Value: "DisplayToken"},
							},
						})
					})
				}
			}()
		case "g45_retrieve_token":
			go func() {
				txtChan := prompt_modal.Instance.Open("", lang.Translate("Enter amount"), key.HintNumeric)
				for txt := range txtChan {
					amount := utils.ShiftNumber{Decimals: int(p.token.Decimals)}
					err := amount.Parse(txt)
					if err != nil {
						return
					}

					scId := p.token.GetHash()
					build_tx_modal.Instance.OpenWithRandomAddr(scId, func(randomAddr string, open func(txPayload build_tx_modal.TxPayload)) {
						open(build_tx_modal.TxPayload{
							Ringsize: 2,
							SCArgs: rpc.Arguments{
								{Name: rpc.SCACTION, DataType: rpc.DataUint64, Value: uint64(rpc.SC_CALL)},
								{Name: rpc.SCID, DataType: rpc.DataHash, Value: scId},
								{Name: "entrypoint", DataType: rpc.DataString, Value: "RetrieveToken"},
								{Name: "amount", DataType: rpc.DataUint64, Value: amount.Number},
							},
						})
					})
				}
			}()
		}

		if err != nil {
			notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else if successMsg != "" {
			notification_modals.SuccessInstance.SetText(lang.Translate("Success"), successMsg)
			notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
			p.tokenMenuSelect.SelectModal.Modal.SetVisible(false)
		}
	}

	if p.sendReceiveButtons.ButtonSend.Clicked() {
		page_instance.pageSendForm.SetToken(p.token)
		page_instance.pageSendForm.ClearForm()
		page_instance.pageRouter.SetCurrent(PAGE_SEND_FORM)
		page_instance.header.AddHistory(PAGE_SEND_FORM)
	}

	if p.sendReceiveButtons.ButtonReceive.Clicked() {
		page_instance.pageRouter.SetCurrent(PAGE_RECEIVE_FORM)
		page_instance.header.AddHistory(PAGE_RECEIVE_FORM)
	}

	widgets := []layout.Widget{}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				r := op.Record(gtx.Ops)
				dims := layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							p.tokenImage.Image.Src = p.token.LoadImageOp()
							gtx.Constraints.Max.X = gtx.Dp(50)
							gtx.Constraints.Max.Y = gtx.Dp(50)
							return p.tokenImage.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							editor := material.Editor(th, p.scIdEditor, "")
							editor.TextSize = unit.Sp(14)
							return editor.Layout(gtx)
						}),
					)
				})
				c := r.Stop()

				paint.FillShape(gtx.Ops, theme.Current.ListBgColor,
					clip.UniformRRect(
						image.Rectangle{Max: dims.Size},
						gtx.Dp(15),
					).Op(gtx.Ops))

				c.Add(gtx.Ops)
				return dims
			}),
		)
	})

	switch p.token.StandardType {
	case sc.G45_AT_TYPE, sc.G45_FAT_TYPE, sc.G45_NFT_TYPE:
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return p.g45DisplayContainer.Layout(gtx, th)
		})
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.balanceContainer.Layout(gtx, th)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.sendReceiveButtons.Layout(gtx, th)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		text := make(map[string]string)
		text["txs"] = lang.Translate("Transactions")
		text["info"] = lang.Translate("Info")
		p.tabBars.Colors = theme.Current.TabBarsColors
		return p.tabBars.Layout(gtx, th, unit.Sp(18), text)
	})

	if p.tabBars.Key == "txs" {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return p.txBar.Layout(gtx, th)
		})

		if len(p.txItems) == 0 {
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(th, unit.Sp(16), lang.Translate("You don't have any txs. Try adjusting filtering options or wait for wallet to sync."))
				lbl.Color = theme.Current.TextMuteColor
				return lbl.Layout(gtx)
			})
		}

		for i := range p.txItems {
			idx := i
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				return p.txItems[idx].Layout(gtx, th)
			})
		}
	}

	if p.tabBars.Key == "info" {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return p.tokenInfo.Layout(gtx, th)
		})
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
	})

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

type TokenMenuSelect struct {
	SelectModal *prefabs.SelectModal
}

func NewTokenMenuSelect() *TokenMenuSelect {
	var items []*prefabs.SelectListItem

	// smart contract options
	// g45_nft
	showIcon, _ := widget.NewIcon(icons.ActionVisibility)
	items = append(items, prefabs.NewSelectListItem("g45_display_nft", prefabs.ListItemMenuItem{
		Icon:  showIcon,
		Title: "Display NFT", //@lang.Translate("Display NFT")
	}.Layout))

	hideIcon, _ := widget.NewIcon(icons.ActionVisibilityOff)
	items = append(items, prefabs.NewSelectListItem("g45_retrieve_nft", prefabs.ListItemMenuItem{
		Icon:  hideIcon,
		Title: "Retrieve NFT", //@lang.Translate("Retrieve NFT")
	}.Layout))

	// g45_fat, g45_at
	items = append(items, prefabs.NewSelectListItem("g45_display_token", prefabs.ListItemMenuItem{
		Icon:  showIcon,
		Title: "Display Tokens", //@lang.Translate("Display Tokens")
	}.Layout))

	items = append(items, prefabs.NewSelectListItem("g45_retrieve_token", prefabs.ListItemMenuItem{
		Icon:  hideIcon,
		Title: "Retrieve Tokens", //@lang.Translate("Retrieve Tokens")
	}.Layout))

	refreshIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	items = append(items, prefabs.NewSelectListItem("refresh_cache", prefabs.ListItemMenuItem{
		Icon:  refreshIcon,
		Title: "Refresh cache", //@lang.Translate("Refresh cache")
	}.Layout))

	addFavIcon, _ := widget.NewIcon(icons.ToggleStarBorder)
	items = append(items, prefabs.NewSelectListItem("add_favorite", prefabs.ListItemMenuItem{
		Icon:  addFavIcon,
		Title: "Add to favorites", //@lang.Translate("Add to favorites")
	}.Layout))

	delFavIcon, _ := widget.NewIcon(icons.ToggleStar)
	items = append(items, prefabs.NewSelectListItem("remove_favorite", prefabs.ListItemMenuItem{
		Icon:  delFavIcon,
		Title: "Remove from favorites", //@lang.Translate("Remove from favorites")
	}.Layout))

	editIcon, _ := widget.NewIcon(icons.ActionInput)
	items = append(items, prefabs.NewSelectListItem("edit_token", prefabs.ListItemMenuItem{
		Icon:  editIcon,
		Title: "Edit token", //@lang.Translate("Edit token")
	}.Layout))

	deleteIcon, _ := widget.NewIcon(icons.ActionDelete)
	items = append(items, prefabs.NewSelectListItem("remove_token", prefabs.ListItemMenuItem{
		Icon:  deleteIcon,
		Title: "Remove token", //@lang.Translate("Remove token")
	}.Layout))

	selectModal := prefabs.NewSelectModal()
	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			var filteredItems []*prefabs.SelectListItem
			for _, item := range items {
				add := true

				token := page_instance.pageSCToken.token

				isFav := false
				if token != nil && token.IsFavorite.Valid {
					isFav = token.IsFavorite.Bool
				}

				standardType := sc.UNKNOWN_TYPE
				if token != nil {
					standardType = token.StandardType
				}

				switch item.Key {
				case "add_favorite":
					add = !isFav
				case "remove_favorite":
					add = isFav
				case "g45_display_nft", "g45_retrieve_nft":
					add = standardType == sc.G45_NFT_TYPE
				case "g45_display_token", "g45_retrieve_token":
					add = standardType == sc.G45_AT_TYPE || standardType == sc.G45_FAT_TYPE
				}

				if add {
					filteredItems = append(filteredItems, item)
				}
			}

			selectModal.Layout(gtx, th, filteredItems)
		},
	})

	return &TokenMenuSelect{
		SelectModal: selectModal,
	}
}

type TokenInfoList struct {
	token    *wallet_manager.Token
	infoRows []*prefabs.InfoRow
}

func NewTokenInfoList(token *wallet_manager.Token) *TokenInfoList {
	var infoRows []*prefabs.InfoRow
	for i := 0; i < 5; i++ {
		infoRows = append(infoRows, prefabs.NewInfoRow())
	}

	return &TokenInfoList{
		token:    token,
		infoRows: infoRows,
	}
}

func (t *TokenInfoList) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	token := t.token
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.infoRows[0].Layout(gtx, th, lang.Translate("Name"), token.Name)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.infoRows[1].Layout(gtx, th, lang.Translate("Decimals"), fmt.Sprint(token.Decimals))
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.infoRows[2].Layout(gtx, th, lang.Translate("Symbol"), token.Symbol.String)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			txt := "?"
			if token.MaxSupply.Valid {
				maxSupply := utils.ShiftNumber{Number: uint64(token.MaxSupply.Int64), Decimals: int(token.Decimals)}
				txt = maxSupply.Format()
			}

			return t.infoRows[3].Layout(gtx, th, lang.Translate("Max Supply"), txt)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.infoRows[4].Layout(gtx, th, lang.Translate("SC Standard"), fmt.Sprint(token.StandardType))
		}),
	)
}

type BalanceContainer struct {
	token             *wallet_manager.Token
	balanceEditor     *widget.Editor
	buttonHideBalance *ButtonHideBalance
}

func NewBalanceContainer() *BalanceContainer {
	buttonHideBalance := NewButtonHideBalance()
	balanceEditor := new(widget.Editor)
	balanceEditor.ReadOnly = true
	balanceEditor.SingleLine = true

	return &BalanceContainer{
		buttonHideBalance: buttonHideBalance,
		balanceEditor:     balanceEditor,
	}
}

func (b *BalanceContainer) SetToken(token *wallet_manager.Token) {
	b.token = token
}

func (b *BalanceContainer) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			r := op.Record(gtx.Ops)
			dims := layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis: layout.Vertical,
						}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								lbl := material.Label(th, unit.Sp(14), lang.Translate("Available Balance"))
								lbl.Color = theme.Current.TextMuteColor
								return lbl.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								wallet := wallet_manager.OpenedWallet
								balance, _ := wallet.Memory.Get_Balance_scid(b.token.GetHash())
								amount := utils.ShiftNumber{Number: balance, Decimals: int(b.token.Decimals)}.Format()

								if b.balanceEditor.Text() != amount {
									b.balanceEditor.SetText(amount)
								}

								r := op.Record(gtx.Ops)
								balanceEditor := material.Editor(th, b.balanceEditor, "")
								balanceEditor.TextSize = unit.Sp(34)
								balanceEditor.Font.Weight = font.Bold

								dims := balanceEditor.Layout(gtx)
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
						)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.Y = gtx.Dp(30)
						gtx.Constraints.Min.X = gtx.Dp(30)
						b.buttonHideBalance.Button.Style.Colors = theme.Current.ButtonIconPrimaryColors
						return b.buttonHideBalance.Layout(gtx, th)
					}),
				)
			})
			c := r.Stop()

			paint.FillShape(gtx.Ops, theme.Current.ListBgColor,
				clip.UniformRRect(
					image.Rectangle{Max: dims.Size},
					gtx.Dp(15),
				).Op(gtx.Ops))

			c.Add(gtx.Ops)
			return dims
		}),
	)
}

type G45DisplayContainer struct {
	token *wallet_manager.Token

	ownerEditor  *widget.Editor
	amountEditor *widget.Editor
}

func NewG45DisplayContainer() *G45DisplayContainer {
	ownerEditor := new(widget.Editor)
	ownerEditor.ReadOnly = true

	amountEditor := new(widget.Editor)
	amountEditor.ReadOnly = true

	return &G45DisplayContainer{
		ownerEditor:  ownerEditor,
		amountEditor: amountEditor,
	}
}

func (d *G45DisplayContainer) SetToken(token *wallet_manager.Token) {
	d.token = token
}

func (d *G45DisplayContainer) Load() {
	switch d.token.StandardType {
	case sc.G45_NFT_TYPE:
		d.ownerEditor.SetText("")

		var result rpc.GetSC_Result
		err := walletapi.RPC_Client.RPC.CallResult(context.Background(), "DERO.GetSC", rpc.GetSC_Params{
			SCID:       d.token.SCID,
			Code:       false,
			Variables:  false,
			KeysString: []string{"owner"},
		}, &result)
		if err != nil {
			d.ownerEditor.SetText("--")
			return
		}

		owner, err := utils.DecodeString(result.ValuesString[0])
		if err != nil {
			d.ownerEditor.SetText("--")
			return
		}

		d.ownerEditor.SetText(owner)
	case sc.G45_AT_TYPE, sc.G45_FAT_TYPE:
		d.amountEditor.SetText("")

		wallet := wallet_manager.OpenedWallet
		addr := wallet.Memory.GetAddress().String()
		key := fmt.Sprintf("owner_%s", addr)

		var result rpc.GetSC_Result
		err := walletapi.RPC_Client.RPC.CallResult(context.Background(), "DERO.GetSC", rpc.GetSC_Params{
			SCID:       d.token.SCID,
			Code:       false,
			Variables:  false,
			KeysString: []string{key},
		}, &result)
		if err != nil {
			d.amountEditor.SetText("--")
			return
		}

		amountDisplayed, err := strconv.ParseUint(result.ValuesString[0], 10, 64)
		if err != nil {
			d.amountEditor.SetText("--")
			return
		}

		amount := utils.ShiftNumber{Number: amountDisplayed, Decimals: int(d.token.Decimals)}
		d.amountEditor.SetText(amount.Format())
	}
}

func (d *G45DisplayContainer) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if d.token == nil { //|| d.amountDisplayed == 0 {
		return layout.Dimensions{}
	}

	r := op.Record(gtx.Ops)
	dims := layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if d.token.StandardType == sc.G45_NFT_TYPE {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), lang.Translate("Owner"))
					lbl.Color = theme.Current.TextMuteColor
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if d.ownerEditor.Text() == "" {
						d.ownerEditor.SetText(lang.Translate("unknown"))
					}

					editor := material.Editor(th, d.ownerEditor, "")
					editor.Font.Weight = font.Bold
					editor.TextSize = unit.Sp(16)
					return editor.Layout(gtx)
				}),
			)
		}

		if d.token.StandardType == sc.G45_AT_TYPE ||
			d.token.StandardType == sc.G45_FAT_TYPE {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), lang.Translate("Amount Displayed"))
					lbl.Color = theme.Current.TextMuteColor
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(th, d.amountEditor, "")
					editor.Font.Weight = font.Bold
					editor.TextSize = unit.Sp(20)
					return editor.Layout(gtx)
				}),
			)
		}

		return layout.Dimensions{}
	})
	c := r.Stop()

	paint.FillShape(gtx.Ops, theme.Current.ListBgColor,
		clip.UniformRRect(
			image.Rectangle{Max: dims.Size},
			gtx.Dp(15),
		).Op(gtx.Ops))

	c.Add(gtx.Ops)
	return dims
}
