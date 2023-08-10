package page_wallet

import (
	"database/sql"
	"fmt"
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
	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/image_modal"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
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

type PageSCToken struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonOpenMenu     *components.Button
	tokenMenuSelect    *TokenMenuSelect
	sendReceiveButtons *SendReceiveButtons
	confirmRemoveToken *prefabs.Confirm
	buttonHideBalance  *ButtonHideBalance
	tabBars            *components.TabBars
	txBar              *TxBar
	getTransfersParams wallet_manager.GetTransfersParams
	txItems            []*TxListItem
	tokenInfo          *TokenInfoList
	balanceEditor      *widget.Editor

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
		Animation: components.NewButtonAnimationScale(.98),
	})

	list := new(widget.List)
	list.Axis = layout.Vertical

	image := prefabs.NewImageHoverClick()

	scIdEditor := new(widget.Editor)
	scIdEditor.WrapPolicy = text.WrapGraphemes
	scIdEditor.ReadOnly = true

	sendReceiveButtons := NewSendReceiveButtons()
	buttonHideBalance := NewButtonHideBalance()
	balanceEditor := new(widget.Editor)
	balanceEditor.ReadOnly = true
	balanceEditor.SingleLine = true

	tabBarsItems := []*components.TabBarsItem{
		components.NewTabBarItem("txs"),
		components.NewTabBarItem("info"),
	}

	tabBars := components.NewTabBars("txs", tabBarsItems)
	txBar := NewTxBar()

	confirmRemoveToken := prefabs.NewConfirm(layout.Center)
	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 2,
		Layout: func(gtx layout.Context, th *material.Theme) {
			confirmRemoveToken.Layout(gtx, th, prefabs.ConfirmText{
				Prompt: lang.Translate("Are you sure?"),
				No:     lang.Translate("NO"),
				Yes:    lang.Translate("YES"),
			})
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
		balanceEditor:      balanceEditor,

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

	p.tokenInfo = NewTokenInfoList(p.token)

	img, err := p.token.LoadImage()
	if err == nil {
		p.tokenImage.Image.Src = paint.NewImageOp(img)
	} else {
		p.tokenImage.Image.Src = theme.Current.TokenImage
	}

	page_instance.header.Title = func() string { return p.token.Name }
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

	for _, entry := range entries {
		txItems = append(txItems, NewTxListItem(entry))
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

	if p.tokenImage.Clickable.Clicked() {
		image_modal.Instance.Open(p.token.Name, p.tokenImage.Image.Src)
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
		page_instance.pageSendForm.clearForm()
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
									lbl.Color = theme.Current.TextMuteColor
									return lbl.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									balance := utils.ShiftNumber{Number: p.balance, Decimals: int(p.token.Decimals)}.Format()

									balanceEditor := material.Editor(th, p.balanceEditor, "")
									balanceEditor.TextSize = unit.Sp(34)
									balanceEditor.Font.Weight = font.Bold

									if balanceEditor.Editor.Text() != balance {
										balanceEditor.Editor.SetText(balance)
									}

									dims := balanceEditor.Layout(gtx)

									if settings.App.HideBalance {
										paint.FillShape(gtx.Ops, theme.Current.HideBalanceBgColor, clip.Rect{
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
							p.buttonHideBalance.Button.Style.Colors = theme.Current.ButtonIconPrimaryColors
							return p.buttonHideBalance.Layout(gtx, th)
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
				lbl := material.Label(th, unit.Sp(16), lang.Translate("You don't have any txs. Try adjusting filering options or wait for wallet to sync."))
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

type TokenInfoList struct {
	nameEditor         *widget.Editor
	decimalsEditor     *widget.Editor
	symbolEditor       *widget.Editor
	maxSupplyEditor    *widget.Editor
	standardTypeEditor *widget.Editor
}

func NewTokenInfoList(token *wallet_manager.Token) *TokenInfoList {
	nameEditor := &widget.Editor{ReadOnly: true}
	nameEditor.SetText(token.Name)

	decimalsEditor := &widget.Editor{ReadOnly: true}
	decimalsEditor.SetText(fmt.Sprint(token.Decimals))

	symbolEditor := &widget.Editor{ReadOnly: true}
	symbolEditor.SetText(fmt.Sprint(token.Symbol.String))

	maxSupplyEditor := &widget.Editor{ReadOnly: true}
	if token.MaxSupply.Valid {
		maxSupply := utils.ShiftNumber{Number: uint64(token.MaxSupply.Int64), Decimals: int(token.Decimals)}
		maxSupplyEditor.SetText(maxSupply.Format())
	} else {
		maxSupplyEditor.SetText("?")
	}

	standardTypeEditor := &widget.Editor{ReadOnly: true}
	standardTypeEditor.SetText(fmt.Sprint(token.StandardType))

	return &TokenInfoList{
		decimalsEditor:     decimalsEditor,
		nameEditor:         nameEditor,
		symbolEditor:       symbolEditor,
		maxSupplyEditor:    maxSupplyEditor,
		standardTypeEditor: standardTypeEditor,
	}
}

func (t *TokenInfoList) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	var flexChilds []layout.FlexChild

	flexChilds = append(flexChilds,
		InfoRowLayout{Editor: t.nameEditor}.Layout(gtx, th, lang.Translate("Name")),
		InfoRowLayout{Editor: t.decimalsEditor}.Layout(gtx, th, lang.Translate("Decimals")),
		InfoRowLayout{Editor: t.symbolEditor}.Layout(gtx, th, lang.Translate("Symbol")),
		InfoRowLayout{Editor: t.maxSupplyEditor}.Layout(gtx, th, lang.Translate("Max Supply")),
		InfoRowLayout{Editor: t.standardTypeEditor}.Layout(gtx, th, lang.Translate("SC Standard")),
	)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, flexChilds...)
}
