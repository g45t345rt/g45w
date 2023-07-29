package page_wallet_select

import (
	"fmt"
	"image"
	"image/color"
	"sort"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
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

type PageSelectWallet struct {
	isActive  bool
	clickable *widget.Clickable

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonWalletCreate *components.Button
	walletList         *WalletList

	modalWalletPassword        *prefabs.PasswordModal
	modalCreateWalletSelection *CreateWalletSelectionModal

	currentWallet *wallet_manager.WalletInfo
}

var _ router.Page = &PageSelectWallet{}

func NewPageSelectWallet() *PageSelectWallet {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	walletList := NewWalletList()

	modalWalletPassword := prefabs.NewPasswordModal()
	modalCreateWalletSelection := NewCreateWalletSelectionModal()

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			modalWalletPassword.Layout(gtx, th)
			modalCreateWalletSelection.Layout(gtx, th)
		},
	})

	addIcon, _ := widget.NewIcon(icons.ContentAddCircleOutline)
	buttonWalletCreate := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            addIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonWalletCreate.Label.Alignment = text.Middle
	buttonWalletCreate.Style.Font.Weight = font.Bold

	return &PageSelectWallet{
		clickable: new(widget.Clickable),

		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonWalletCreate: buttonWalletCreate,
		walletList:         walletList,

		modalWalletPassword:        modalWalletPassword,
		modalCreateWalletSelection: modalCreateWalletSelection,
	}
}

func (p *PageSelectWallet) IsActive() bool {
	return p.isActive
}

func (p *PageSelectWallet) Enter() {
	p.isActive = true
	page_instance.header.SetTitle(lang.Translate("Select wallet"))

	if !page_instance.header.IsHistory(PAGE_SELECT_WALLET) {
		p.animationLeave.Reset()
		p.animationEnter.Start()
	}

	p.Load()
}

func (p *PageSelectWallet) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *PageSelectWallet) Load() {
	wallets := wallet_manager.Wallets
	items := make([]WalletListItem, 0)
	for _, wallet := range wallets {
		items = append(items,
			NewWalletListItem(wallet),
		)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].wallet.Timestamp < items[j].wallet.Timestamp
	})

	p.walletList.items = items
}

func (p *PageSelectWallet) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top: unit.Dp(0), Bottom: unit.Dp(30),
				Left: unit.Dp(30), Right: unit.Dp(30),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						if len(p.walletList.items) == 0 {
							labelNoWallet := material.Label(th, unit.Sp(16), lang.Translate("You didn't add a wallet yet.\nClick 'New Wallet' button to continue."))
							return labelNoWallet.Layout(gtx)
						} else {
							for _, item := range p.walletList.items {
								if item.Clickable.Clicked() {
									p.currentWallet = item.wallet
									p.modalWalletPassword.Modal.SetVisible(true)
								}
							}

							return p.walletList.Layout(gtx, th)
						}
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(30)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if p.buttonWalletCreate.Clicked() {
							p.modalCreateWalletSelection.modal.SetVisible(true)
						}

						p.buttonWalletCreate.Text = lang.Translate("NEW WALLET")
						return p.buttonWalletCreate.Layout(gtx, th)
					}),
				)
			})
		}),
	)

	{
		submitted, text := p.modalWalletPassword.Input.Submitted()
		if submitted {
			err := wallet_manager.OpenWallet(p.currentWallet.Addr, text)
			if err == nil {
				wallet := wallet_manager.OpenedWallet
				wallet.Memory.SetOnlineMode()
				p.modalWalletPassword.Modal.SetVisible(false)
				app_instance.Router.SetCurrent(app_instance.PAGE_WALLET)
			} else {
				if err.Error() == "Invalid Password" {
					p.modalWalletPassword.StartWrongPassAnimation()
				} else {
					//p.modalWalletPassword.Modal.SetVisible(false)
					notification_modals.ErrorInstance.SetText("Error", err.Error())
					notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
				}
			}
		}
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

type CreateWalletSelectionModal struct {
	modal *components.Modal
	list  *widget.List
	items []*CreateWalletListItem
}

func NewCreateWalletSelectionModal() *CreateWalletSelectionModal {
	modal := components.NewModal(components.ModalStyle{
		CloseOnOutsideClick: true,
		CloseOnInsideClick:  false,
		Direction:           layout.S,
		BgColor:             color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Rounded:             components.UniformRounded(unit.Dp(10)),
		Inset:               layout.UniformInset(25),
		Animation:           components.NewModalAnimationUp(),
		Backdrop:            components.NewModalBackground(),
	})

	list := new(widget.List)
	list.Axis = layout.Vertical

	fastIcon, _ := widget.NewIcon(icons.ImageFlashOn)
	newIcon, _ := widget.NewIcon(icons.ContentAddCircle)
	diskIcon, _ := widget.NewIcon(icons.FileFolder)
	seedIcon, _ := widget.NewIcon(icons.EditorShortText)

	items := []*CreateWalletListItem{
		NewCreateWalletListItem("Fast registration", fastIcon, PAGE_CREATE_WALLET_FASTREG_FORM),     //@lang.Translate("Fast registration")
		NewCreateWalletListItem("Create new wallet", newIcon, PAGE_CREATE_WALLET_FORM),              //@lang.Translate("Create new wallet")
		NewCreateWalletListItem("Recover from Disk", diskIcon, PAGE_CREATE_WALLET_DISK_FORM),        //@lang.Translate("Recover from Disk")
		NewCreateWalletListItem("Recover from Seed", seedIcon, PAGE_CREATE_WALLET_SEED_FORM),        //@lang.Translate("Recover from Seed")
		NewCreateWalletListItem("Recover from Hex Seed", seedIcon, PAGE_CREATE_WALLET_HEXSEED_FORM), //@lang.Translate("Recover from Hex Seed")
	}

	return &CreateWalletSelectionModal{
		modal: modal,
		list:  list,
		items: items,
	}
}

func (c *CreateWalletSelectionModal) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return c.modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(10), Bottom: unit.Dp(10),
			Left: unit.Dp(10), Right: unit.Dp(0),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			listStyle := material.List(th, c.list)

			return listStyle.Layout(gtx, len(c.items), func(gtx layout.Context, index int) layout.Dimensions {
				if c.items[index].clickable.Clicked() {
					c.modal.SetVisible(false)

					tag := c.items[index].routerTag
					page_instance.pageRouter.SetCurrent(tag)
					page_instance.header.AddHistory(tag)

					op.InvalidateOp{}.Add(gtx.Ops) // make sure to invalidate if we are closing modal and changing page
				}

				return c.items[index].Layout(gtx, th)
			})
		})
	})
}

type CreateWalletListItem struct {
	text      string
	routerTag string
	icon      *widget.Icon
	clickable *widget.Clickable
}

func NewCreateWalletListItem(text string, icon *widget.Icon, routerTag string) *CreateWalletListItem {
	return &CreateWalletListItem{
		text:      text,
		icon:      icon,
		routerTag: routerTag,
		clickable: new(widget.Clickable),
	}
}

func (c *CreateWalletListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	dims := c.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return c.icon.Layout(gtx, color.NRGBA{A: 255})
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(20), lang.Translate(c.text))
					return lbl.Layout(gtx)
				}),
			)
		})
	})

	if c.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)

		paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 100},
			clip.UniformRRect(
				image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)},
				gtx.Dp(15),
			).Op(gtx.Ops),
		)
	}

	return dims
}

type WalletList struct {
	list  *widget.List
	items []WalletListItem
}

func NewWalletList() *WalletList {
	list := new(widget.List)
	list.Axis = layout.Vertical

	return &WalletList{
		list:  list,
		items: []WalletListItem{},
	}
}

func (w *WalletList) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		clip.UniformRRect(
			image.Rectangle{Max: gtx.Constraints.Max},
			gtx.Dp(unit.Dp(10)),
		).Op(gtx.Ops),
	)

	return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		listStyle := material.List(th, w.list)
		listStyle.AnchorStrategy = material.Overlay
		listStyle.Indicator.MinorWidth = unit.Dp(10)
		listStyle.Indicator.CornerRadius = unit.Dp(5)
		black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
		listStyle.Indicator.Color = black

		return listStyle.Layout(gtx, len(w.items), func(gtx layout.Context, i int) layout.Dimensions {
			return w.items[i].Layout(gtx, th)
		})
	})
}

type WalletListItem struct {
	wallet    *wallet_manager.WalletInfo
	Clickable *widget.Clickable
	rounded   unit.Dp
}

func NewWalletListItem(wallet *wallet_manager.WalletInfo) WalletListItem {
	return WalletListItem{
		wallet:    wallet,
		Clickable: &widget.Clickable{},
		rounded:   unit.Dp(12),
	}
}

func (item *WalletListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return item.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		dims := layout.UniformInset(item.rounded).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Start}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							name := fmt.Sprintf("%s [%s]", lang.Translate("Wallet"), item.wallet.Name)
							lbl := material.Label(th, unit.Sp(18), name)
							lbl.Font.Weight = font.Bold
							return lbl.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							addr := utils.ReduceAddr(item.wallet.Addr)
							lbl := material.Label(th, unit.Sp(15), addr)
							lbl.Color.A = 200
							return lbl.Layout(gtx)
						}),
					)
				}),
			)
		})

		if item.Clickable.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 100},
				clip.UniformRRect(
					image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)},
					gtx.Dp(item.rounded),
				).Op(gtx.Ops),
			)
		}

		return dims
	})
}
