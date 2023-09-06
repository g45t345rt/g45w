package page_wallet_select

import (
	"fmt"
	"image"
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
	"github.com/g45t345rt/g45w/containers/confirm_modal"
	"github.com/g45t345rt/g45w/containers/listselect_modal"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/containers/password_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/pages"
	page_wallet "github.com/g45t345rt/g45w/pages/wallet"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
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

	addIcon, _ := widget.NewIcon(icons.ContentAddCircleOutline)
	buttonWalletCreate := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      addIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonWalletCreate.Label.Alignment = text.Middle
	buttonWalletCreate.Style.Font.Weight = font.Bold

	return &PageSelectWallet{
		clickable: new(widget.Clickable),

		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonWalletCreate: buttonWalletCreate,
		walletList:         walletList,
	}
}

func (p *PageSelectWallet) IsActive() bool {
	return p.isActive
}

func (p *PageSelectWallet) Enter() {
	p.isActive = true
	page_instance.header.Title = func() string { return lang.Translate("Select wallet") }

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
	items := make([]WalletListItem, 0)
	for _, walletInfo := range wallet_manager.Wallets {
		items = append(items,
			NewWalletListItem(walletInfo),
		)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].wallet.Timestamp < items[j].wallet.Timestamp
	})

	for addr, err := range wallet_manager.WalletsErr {
		items = append(items,
			NewWalletListItemErr(err, addr),
		)
	}

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
									if item.wallet != nil {
										p.currentWallet = item.wallet
										password_modal.Instance.SetVisible(true)
									} else {
										go func() {
											yesChan := confirm_modal.Instance.Open(confirm_modal.ConfirmText{
												Prompt: lang.Translate("Delete wallet?"),
											})

											for yes := range yesChan {
												if yes {
													err := wallet_manager.DeleteWallet(item.addr)
													if err == nil {
														notification_modals.ErrorInstance.SetText("Error", err.Error())
														notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
														p.Load()
													}
												}
											}
										}()
									}
								}
							}

							return p.walletList.Layout(gtx, th)
						}
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(30)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if p.buttonWalletCreate.Clicked() {
							go func() {
								fastIcon, _ := widget.NewIcon(icons.ImageFlashOn)
								newIcon, _ := widget.NewIcon(icons.ContentAddCircle)
								diskIcon, _ := widget.NewIcon(icons.FileFolder)
								seedIcon, _ := widget.NewIcon(icons.EditorShortText)

								keyChan := listselect_modal.Instance.Open([]*listselect_modal.SelectListItem{
									listselect_modal.NewSelectListItem(PAGE_CREATE_WALLET_FASTREG_FORM,
										listselect_modal.NewItemText(fastIcon, lang.Translate("Fast registration")).Layout,
									),
									listselect_modal.NewSelectListItem(PAGE_CREATE_WALLET_FORM,
										listselect_modal.NewItemText(newIcon, lang.Translate("Create new wallet")).Layout,
									),
									listselect_modal.NewSelectListItem(PAGE_CREATE_WALLET_DISK_FORM,
										listselect_modal.NewItemText(diskIcon, lang.Translate("Recover from Disk")).Layout,
									),
									listselect_modal.NewSelectListItem(PAGE_CREATE_WALLET_SEED_FORM,
										listselect_modal.NewItemText(seedIcon, lang.Translate("Recover from Seed")).Layout,
									),
									listselect_modal.NewSelectListItem(PAGE_CREATE_WALLET_HEXSEED_FORM,
										listselect_modal.NewItemText(seedIcon, lang.Translate("Recover from Hex Seed")).Layout,
									),
								})

								for key := range keyChan {
									page_instance.pageRouter.SetCurrent(key)
									page_instance.header.AddHistory(key)
								}
							}()
						}

						p.buttonWalletCreate.Text = lang.Translate("NEW WALLET")
						p.buttonWalletCreate.Style.Colors = theme.Current.ButtonPrimaryColors
						return p.buttonWalletCreate.Layout(gtx, th)
					}),
				)
			})
		}),
	)

	{
		submitted, text := password_modal.Instance.Input.Submitted()
		if submitted {
			go func() {
				password_modal.Instance.SetLoading(true)
				err := wallet_manager.OpenWallet(p.currentWallet.Addr, text)
				password_modal.Instance.SetLoading(false)
				if err == nil {
					wallet := wallet_manager.OpenedWallet
					wallet.Memory.SetOnlineMode()
					password_modal.Instance.SetVisible(false)
					// important reset wallet pages to initial state
					app_instance.Router.Pages[pages.PAGE_WALLET] = page_wallet.New()
					app_instance.Router.SetCurrent(pages.PAGE_WALLET)
				} else {
					if err.Error() == "Invalid Password" {
						password_modal.Instance.StartWrongPassAnimation()
					} else {
						//p.modalWalletPassword.Modal.SetVisible(false)
						notification_modals.ErrorInstance.SetText("Error", err.Error())
						notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
					}
				}
			}()
		}
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
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
	paint.FillShape(gtx.Ops, theme.Current.ListBgColor,
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
		listStyle.Indicator.Color = theme.Current.ListScrollBarBgColor

		return listStyle.Layout(gtx, len(w.items), func(gtx layout.Context, i int) layout.Dimensions {
			return w.items[i].Layout(gtx, th)
		})
	})
}

type WalletListItem struct {
	wallet    *wallet_manager.WalletInfo
	Clickable *widget.Clickable
	rounded   unit.Dp
	err       error
	addr      string
}

func NewWalletListItem(wallet *wallet_manager.WalletInfo) WalletListItem {
	return WalletListItem{
		wallet:    wallet,
		Clickable: &widget.Clickable{},
		rounded:   unit.Dp(12),
	}
}

func NewWalletListItemErr(err error, addr string) WalletListItem {
	return WalletListItem{
		Clickable: &widget.Clickable{},
		rounded:   unit.Dp(12),
		err:       err,
		addr:      addr,
	}
}

func (item *WalletListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return item.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		dims := layout.UniformInset(item.rounded).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Start}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					if item.err != nil {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								addr := utils.ReduceAddr(item.addr)
								lbl := material.Label(th, unit.Sp(18), addr)
								lbl.Font.Weight = font.Bold
								return lbl.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								lbl := material.Label(th, unit.Sp(11), item.err.Error())
								lbl.Color = theme.Current.TextMuteColor
								return lbl.Layout(gtx)
							}),
						)
					}

					if item.wallet != nil {
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
								lbl.Color = theme.Current.TextMuteColor
								return lbl.Layout(gtx)
							}),
						)
					}

					return layout.Dimensions{}
				}),
			)
		})

		if item.Clickable.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
			paint.FillShape(gtx.Ops, theme.Current.ListItemHoverBgColor,
				clip.UniformRRect(
					image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)},
					gtx.Dp(item.rounded),
				).Op(gtx.Ops),
			)
		}

		return dims
	})
}
