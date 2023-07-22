package page_wallet

import (
	"database/sql"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
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

	buttonOpenMenu  *components.Button
	tokenMenuSelect *TokenMenuSelect

	token *wallet_manager.Token

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

	return &PageSCToken{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonOpenMenu:  buttonOpenMenu,
		tokenMenuSelect: NewTokenMenuSelect(),

		list: list,
	}
}

func (p *PageSCToken) IsActive() bool {
	return p.isActive
}

func (p *PageSCToken) Enter() {
	p.isActive = true

	page_instance.header.SetTitle(lang.Translate("Token"))
	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		scId := utils.ReduceTxId(p.token.SCID)
		lbl := material.Label(th, unit.Sp(16), scId)
		return lbl.Layout(gtx)
	}
	page_instance.header.ButtonRight = p.buttonOpenMenu

	if !page_instance.header.IsHistory(PAGE_SC_TOKEN) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageSCToken) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
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
			err = wallet.DelToken(p.token.ID)
			successMsg = lang.Translate("Token removed.")
			page_instance.header.GoBack()
		}

		if err != nil {
			notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {
			notification_modals.SuccessInstance.SetText(lang.Translate("Success"), successMsg)
			notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
			p.tokenMenuSelect.SelectModal.Modal.SetVisible(false)
		}
	}

	widgets := []layout.Widget{}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				r := op.Record(gtx.Ops)
				dims := layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Label(th, unit.Sp(22), p.token.Name)
							label.Color = color.NRGBA{A: 255}
							return label.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Label(th, unit.Sp(16), p.token.SCID)
							label.Color = color.NRGBA{A: 150}
							return label.Layout(gtx)
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
