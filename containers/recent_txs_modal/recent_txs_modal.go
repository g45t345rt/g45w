package recent_txs_modal

import (
	"fmt"
	"image"
	"strings"
	"time"

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
	"gioui.org/x/browser"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/confirm_modal"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type RecentTxsModal struct {
	modal       *components.Modal
	list        *widget.List
	buttonClear *components.Button

	txItems []TxItem
}

var Instance *RecentTxsModal

func LoadInstance() {
	modal := components.NewModal(components.ModalStyle{
		CloseOnOutsideClick: true,
		CloseOnInsideClick:  false,
		Direction:           layout.N,
		Rounded: components.Rounded{
			SW: unit.Dp(10), SE: unit.Dp(10),
		},
		Animation: components.NewModalAnimationDown(),
	})

	list := new(widget.List)
	list.Axis = layout.Vertical

	deleteIcon, _ := widget.NewIcon(icons.ContentDeleteSweep)
	buttonClear := components.NewButton(components.ButtonStyle{
		Icon:      deleteIcon,
		Animation: components.NewButtonAnimationScale(.95),
	})

	Instance = &RecentTxsModal{
		modal:       modal,
		list:        list,
		buttonClear: buttonClear,
	}

	Instance.startCheckingPendingTxs()

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 2,
		Layout: func(gtx layout.Context, th *material.Theme) {
			Instance.layout(gtx, th)
		},
	})
}

func (r *RecentTxsModal) startCheckingPendingTxs() {
	w := app_instance.Window
	ticker := time.NewTicker(15 * time.Second)

	go func() {
		for range ticker.C {
			wallet := wallet_manager.OpenedWallet
			if wallet != nil {
				updated, err := wallet.UpdatePendingOutgoingTxs()
				if err != nil {
					fmt.Println(err)
				}

				if updated > 0 && r.modal.Visible {
					r.LoadOutgoingTxs()
					w.Invalidate()
				}
			}
		}
	}()
}

func (r *RecentTxsModal) LoadOutgoingTxs() error {
	r.txItems = make([]TxItem, 0)

	wallet := wallet_manager.OpenedWallet
	if wallet != nil {
		outgoingTxs, err := wallet.GetOutgoingTxs(wallet_manager.GetOutgoingTxsParams{
			OrderBy:    "timestamp",
			Descending: true,
			Limit:      uint64(10),
		})
		if err != nil {
			return err
		}

		for _, tx := range outgoingTxs {
			r.txItems = append(r.txItems, *NewTxItem(tx))
		}
	}

	return nil
}

func (r *RecentTxsModal) SetVisible(visible bool) {
	if visible {
		r.LoadOutgoingTxs()
	}

	r.modal.SetVisible(visible)
}

func (r *RecentTxsModal) layout(gtx layout.Context, th *material.Theme) {
	wallet := wallet_manager.OpenedWallet

	if r.buttonClear.Clicked(gtx) {
		go func() {
			yesChan := confirm_modal.Instance.Open(confirm_modal.ConfirmText{
				Prompt: lang.Translate("Are you sure you want to clear outgoing txs?"),
			})

			if <-yesChan {
				err := wallet.ClearOutgoingTxs()
				if err != nil {
					notification_modal.Open(notification_modal.Params{
						Type:  notification_modal.ERROR,
						Title: lang.Translate("Error"),
						Text:  err.Error(),
					})
				} else {
					notification_modal.Open(notification_modal.Params{
						Type:       notification_modal.SUCCESS,
						Title:      lang.Translate("Success"),
						Text:       lang.Translate("Outgoing txs cleared."),
						CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
					})
					r.LoadOutgoingTxs()
				}
			}
		}()
	}

	r.buttonClear.Disabled = wallet == nil

	r.modal.Style.Colors = theme.Current.ModalColors
	r.modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: theme.PagePadding, Bottom: theme.PagePadding,
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(20), fmt.Sprintf("%s (%d)", lang.Translate("Outgoing Transactions"), len(r.txItems)))
							lbl.Font.Weight = font.Bold
							return lbl.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							r.buttonClear.Style.Colors = theme.Current.ModalButtonColors
							return r.buttonClear.Layout(gtx, th)
						}),
					)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.Y = gtx.Dp(250)
					if wallet == nil {
						lbl := material.Label(th, unit.Sp(16), lang.Translate("To view transactions, a wallet must be opened."))
						lbl.Color = theme.Current.TextMuteColor
						return lbl.Layout(gtx)
					} else {
						if len(r.txItems) == 0 {
							lbl := material.Label(th, unit.Sp(16), lang.Translate("You don't have any outgoing transactions yet."))
							lbl.Color = theme.Current.TextMuteColor
							return lbl.Layout(gtx)
						}

						listStyle := material.List(th, r.list)
						listStyle.AnchorStrategy = material.Overlay

						return listStyle.Layout(gtx, len(r.txItems), func(gtx layout.Context, index int) layout.Dimensions {
							bottomInset := 0
							if index < len(r.txItems)-1 {
								bottomInset = 5
							}

							return layout.Inset{
								Bottom: unit.Dp(bottomInset), Right: unit.Dp(15),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return r.txItems[index].Layout(gtx, th)
							})
						})
					}
				}),
			)
		})
	})
}

type TxItem struct {
	tx             wallet_manager.OutgoingTx
	buttonOpen     *components.Button
	buttonRemove   *components.Button
	listItemSelect *prefabs.ListItemSelect
	clickable      *widget.Clickable
}

func NewTxItem(tx wallet_manager.OutgoingTx) *TxItem {
	openIcon, _ := widget.NewIcon(icons.ActionOpenInBrowser)

	buttonOpen := components.NewButton(components.ButtonStyle{
		Icon:      openIcon,
		Rounded:   components.UniformRounded(unit.Dp(5)),
		TextSize:  unit.Sp(14),
		Inset:     layout.UniformInset(unit.Dp(5)),
		Animation: components.NewButtonAnimationDefault(),
	})

	remoteIcon, _ := widget.NewIcon(icons.ActionDelete)
	buttonRemove := components.NewButton(components.ButtonStyle{
		Icon:      remoteIcon,
		Rounded:   components.UniformRounded(unit.Dp(5)),
		TextSize:  unit.Sp(14),
		Inset:     layout.UniformInset(unit.Dp(5)),
		Animation: components.NewButtonAnimationDefault(),
	})

	return &TxItem{
		tx:             tx,
		buttonOpen:     buttonOpen,
		buttonRemove:   buttonRemove,
		listItemSelect: prefabs.NewListItemSelect(),
		clickable:      &widget.Clickable{},
	}
}

func (item *TxItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	var status string
	confirmations := uint64(0)
	txId := item.tx.TxId

	switch item.tx.Status.String {
	case "valid":

		height := uint64(walletapi.Get_Daemon_Height())

		if height > 0 {
			confirmations = height - uint64(item.tx.BlockHeight.Int64)
		} else {
			confirmations = 0
		}

		value := lang.Translate("{} confirmations")
		status = strings.Replace(value, "{}", fmt.Sprint(confirmations), -1)
	case "invalid":
		status = lang.Translate("Invalid transaction")
	default:
		status = lang.Translate("Checking transaction...")
	}

	date := time.Unix(item.tx.Timestamp.Int64, 0)

	if item.buttonOpen.Clicked(gtx) {
		go func() {
			url := fmt.Sprintf("https://explorer.dero.io/tx/%s", txId)
			err := browser.OpenUrl(url)
			if err != nil {
				notification_modal.Open(notification_modal.Params{
					Type:  notification_modal.ERROR,
					Title: lang.Translate("Error"),
					Text:  err.Error(),
				})
			}
		}()
	}

	if item.buttonRemove.Clicked(gtx) {
		wallet := wallet_manager.OpenedWallet
		err := wallet.DelOutgoingTx(item.tx.TxId)

		if err != nil {
			notification_modal.Open(notification_modal.Params{
				Type:  notification_modal.ERROR,
				Title: lang.Translate("Error"),
				Text:  err.Error(),
			})
		} else {
			notification_modal.Open(notification_modal.Params{
				Type:       notification_modal.SUCCESS,
				Title:      lang.Translate("Success"),
				Text:       lang.Translate("Transaction ref remove."),
				CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
			})
			Instance.LoadOutgoingTxs()
		}
	}

	if item.clickable.Clicked(gtx) {
		item.listItemSelect.Toggle()
	}

	r := op.Record(gtx.Ops)
	dims := item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(5), Bottom: unit.Dp(5),
			Left: unit.Dp(8), Right: unit.Dp(8),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			var flexChilds []layout.FlexChild

			flexChilds = append(flexChilds, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), utils.ReduceTxId(txId))
						lbl.Font.Weight = font.Bold
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), status)
						lbl.Color = theme.Current.TextMuteColor
						return lbl.Layout(gtx)
					}),
				)
			}))

			flexChilds = append(flexChilds, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), fmt.Sprint(item.tx.BlockHeight.Int64))
						lbl.Alignment = text.End
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), lang.TimeAgo(date))
						lbl.Alignment = text.End
						lbl.Color = theme.Current.TextMuteColor
						return lbl.Layout(gtx)
					}),
				)
			}))

			r := op.Record(gtx.Ops)
			dims := layout.Flex{
				Axis:      layout.Horizontal,
				Spacing:   layout.SpaceBetween,
				Alignment: layout.Middle,
			}.Layout(gtx,
				flexChilds...,
			)
			c := r.Stop()

			c.Add(gtx.Ops)

			item.buttonOpen.Style.Colors = theme.Current.ButtonPrimaryColors
			item.buttonRemove.Style.Colors = theme.Current.ButtonPrimaryColors
			item.listItemSelect.Layout(gtx, th, []*components.Button{item.buttonOpen, item.buttonRemove})

			return dims
		})
	})
	c := r.Stop()

	if item.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
		paint.FillShape(gtx.Ops, theme.Current.ListItemHoverBgColor,
			clip.UniformRRect(
				image.Rectangle{Max: dims.Size},
				gtx.Dp(5)).Op(gtx.Ops),
		)
	}

	c.Add(gtx.Ops)
	return dims
}
