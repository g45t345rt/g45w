package page_wallet

import (
	"database/sql"
	"fmt"
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
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

	buttonOpenMenu     *components.Button
	tokenMenuSelect    *TokenMenuSelect
	sendReceiveButtons *SendReceiveButtons
	confirmRemoveToken *components.Confirm
	buttonHideBalance  *ButtonHideBalance
	tabBars            *components.TabBars
	txBar              *TxBar
	getTransfersParams wallet_manager.GetTransfersParams
	txItems            []*TxListItem

	token      *wallet_manager.Token
	tokenImage *prefabs.ImageHoverClick
	scIdEditor *widget.Editor

	balance uint64

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
		TextColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		Animation: components.NewButtonAnimationScale(.98),
	})

	list := new(widget.List)
	list.Axis = layout.Vertical

	src, _ := assets.GetImage("token.png")
	image := prefabs.NewImageHoverClick(paint.NewImageOp(src))

	scIdEditor := new(widget.Editor)
	scIdEditor.WrapPolicy = text.WrapGraphemes
	scIdEditor.ReadOnly = true

	sendReceiveButtons := NewSendReceiveButtons()
	buttonHideBalance := NewButtonHideBalance()

	tabBarsItems := []*components.TabBarsItem{
		components.NewTabBarItem("txs", components.TabBarItemStyle{
			TextSize: unit.Sp(18),
		}),
	}

	tabBars := components.NewTabBars("txs", tabBarsItems)

	txBar := NewTxBar()

	confirmRemoveToken := components.NewConfirm(layout.Center)
	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 2,
		Layout: func(gtx layout.Context, th *material.Theme) {
			confirmRemoveToken.Prompt = lang.Translate("Are you sure?")
			confirmRemoveToken.NoText = lang.Translate("NO")
			confirmRemoveToken.YesText = lang.Translate("YES")
			confirmRemoveToken.Layout(gtx, th)
		},
	})

	return &PageSCToken{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonOpenMenu:     buttonOpenMenu,
		tokenMenuSelect:    NewTokenMenuSelect(),
		tokenImage:         image,
		scIdEditor:         scIdEditor,
		sendReceiveButtons: sendReceiveButtons,
		confirmRemoveToken: confirmRemoveToken,
		buttonHideBalance:  buttonHideBalance,
		tabBars:            tabBars,
		txBar:              txBar,

		list: list,
	}
}

func (p *PageSCToken) IsActive() bool {
	return p.isActive
}

func (p *PageSCToken) Enter() {
	p.isActive = true

	wallet := wallet_manager.OpenedWallet
	scId := crypto.HashHexToHash(p.token.SCID)
	wallet.Memory.TokenAdd(scId) // we don't check error because the only possible error is if the token was already added

	page_instance.header.SetTitle(p.token.Name)
	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		scId := utils.ReduceTxId(p.token.SCID)
		if p.token.Symbol.Valid {
			scId = fmt.Sprintf("%s (%s)", scId, p.token.Symbol.String)
		}

		lbl := material.Label(th, unit.Sp(16), scId)
		return lbl.Layout(gtx)
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

func (p *PageSCToken) RefreshBalance() {
	wallet := wallet_manager.OpenedWallet
	scId := crypto.HashHexToHash(p.token.SCID)
	p.balance, _ = wallet.Memory.Get_Balance_scid(scId)
}

func (p *PageSCToken) LoadTxs() {
	wallet := wallet_manager.OpenedWallet
	entries := wallet.GetTransfers(p.token.SCID, p.getTransfersParams)

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

	if p.confirmRemoveToken.ClickedYes() {
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

	selected, key := p.tokenMenuSelect.SelectModal.Selected()
	if selected {
		wallet := wallet_manager.OpenedWallet
		var err error
		var successMsg = ""

		switch key {
		case "add_favorite":
			p.token.IsFavorite = sql.NullBool{Bool: true, Valid: true}
			err = wallet.UpdateToken(*p.token)
			successMsg = lang.Translate("Token added to favorites.")
		case "remove_favorite":
			p.token.IsFavorite = sql.NullBool{Bool: false, Valid: true}
			err = wallet.UpdateToken(*p.token)
			successMsg = lang.Translate("Token removed from favorites.")
		case "remove_token":
			p.confirmRemoveToken.SetVisible(true)
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
		page_instance.pageSendForm.token = *p.token
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
							gtx.Constraints.Max.X = gtx.Dp(50)
							gtx.Constraints.Max.Y = gtx.Dp(50)
							return p.tokenImage.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							editor := material.Editor(th, p.scIdEditor, "")
							return editor.Layout(gtx)
						}),
					)
				})
				c := r.Stop()

				paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255},
					clip.UniformRRect(
						image.Rectangle{Max: dims.Size},
						gtx.Dp(15),
					).Op(gtx.Ops))

				c.Add(gtx.Ops)
				return dims
			}),
		)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
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
									lbl.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 150}
									return lbl.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									balance := utils.ShiftNumber{Number: p.balance, Decimals: int(p.token.Decimals)}
									lbl := material.Label(th, unit.Sp(34), balance.Format())
									lbl.Font.Weight = font.Bold

									dims := lbl.Layout(gtx)

									if settings.App.HideBalance {
										paint.FillShape(gtx.Ops, color.NRGBA{R: 200, G: 200, B: 200, A: 255}, clip.Rect{
											Max: dims.Size,
										}.Op())
									}

									return dims
								}),
							)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.Y = gtx.Dp(30)
							gtx.Constraints.Min.X = gtx.Dp(30)
							return p.buttonHideBalance.Layout(gtx, th)
						}),
					)
				})
				c := r.Stop()

				paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255},
					clip.UniformRRect(
						image.Rectangle{Max: dims.Size},
						gtx.Dp(15),
					).Op(gtx.Ops))

				c.Add(gtx.Ops)
				return dims
			}),
		)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.sendReceiveButtons.Layout(gtx, th)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		text := make(map[string]string)
		text["txs"] = lang.Translate("Transactions")
		return p.tabBars.Layout(gtx, th, text)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.txBar.Layout(gtx, th)
	})

	if len(p.txItems) == 0 {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), lang.Translate("You don't have any txs. Try adjusting filering options."))
			return lbl.Layout(gtx)
		})
	}

	for i := range p.txItems {
		idx := i
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return p.txItems[idx].Layout(gtx, th)
		})
	}

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
	addFavIcon, _ := widget.NewIcon(icons.ToggleStarBorder)
	items = append(items, prefabs.NewSelectListItem("add_favorite", FolderMenuItem{
		Icon:  addFavIcon,
		Title: "Add to favorites", //@lang.Translate("Add to favorites")
	}.Layout))

	delFavIcon, _ := widget.NewIcon(icons.ToggleStar)
	items = append(items, prefabs.NewSelectListItem("remove_favorite", FolderMenuItem{
		Icon:  delFavIcon,
		Title: "Remove from favorites", //@lang.Translate("Remove from favorites")
	}.Layout))

	editIcon, _ := widget.NewIcon(icons.ActionInput)
	items = append(items, prefabs.NewSelectListItem("edit_token", FolderMenuItem{
		Icon:  editIcon,
		Title: "Edit token", //@lang.Translate("Edit token")
	}.Layout))

	deleteIcon, _ := widget.NewIcon(icons.ActionDelete)
	items = append(items, prefabs.NewSelectListItem("remove_token", FolderMenuItem{
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

				switch item.Key {
				case "add_favorite":
					add = !isFav
				case "remove_favorite":
					add = isFav
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
