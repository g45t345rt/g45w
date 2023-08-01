package recent_txs_modal

import (
	"fmt"
	"image"
	"image/color"
	"strings"
	"time"

	"gioui.org/font"
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
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type RecentTxsModal struct {
	modal        *components.Modal
	list         *widget.List
	buttonClear  *components.Button
	confirmClear *components.Confirm

	txItems []TxItem
}

var Instance *RecentTxsModal

func LoadInstance() {
	modal := components.NewModal(components.ModalStyle{
		CloseOnOutsideClick: true,
		CloseOnInsideClick:  false,
		Direction:           layout.N,
		BgColor:             color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Rounded: components.Rounded{
			SW: unit.Dp(10), SE: unit.Dp(10),
		},
		Animation: components.NewModalAnimationDown(),
		Backdrop:  components.NewModalBackground(),
	})

	confirmClear := components.NewConfirm(layout.Center)

	list := new(widget.List)
	list.Axis = layout.Vertical

	deleteIcon, _ := widget.NewIcon(icons.ContentDeleteSweep)
	buttonClear := components.NewButton(components.ButtonStyle{
		Icon:           deleteIcon,
		TextColor:      color.NRGBA{R: 0, G: 0, B: 0, A: 100},
		HoverTextColor: &color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		Animation:      components.NewButtonAnimationScale(.95),
	})

	Instance = &RecentTxsModal{
		modal:        modal,
		list:         list,
		buttonClear:  buttonClear,
		confirmClear: confirmClear,
	}

	Instance.startCheckingPendingTxs()

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 2,
		Layout: func(gtx layout.Context, th *material.Theme) {
			Instance.layout(gtx, th)

			confirmClear.Prompt = lang.Translate("Are you sure you want to clear outgoing txs?")
			confirmClear.NoText = lang.Translate("NO")
			confirmClear.YesText = lang.Translate("YES")
			confirmClear.Layout(gtx, th)
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

				if updated > 0 {
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
		limit := uint64(10)
		outgoingTxs, err := wallet.GetOutgoingTxs(wallet_manager.GetOutgoingTxsParams{
			OrderBy:    "timestamp",
			Descending: true,
			Limit:      &limit,
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
	if r.buttonClear.Clicked() {
		r.confirmClear.SetVisible(true)
	}

	wallet := wallet_manager.OpenedWallet
	if r.confirmClear.ClickedYes() {
		err := wallet.ClearOutgoingTxs()
		if err != nil {
			notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {
			notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("Outgoing txs cleared."))
			notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		}
	}

	r.buttonClear.Disabled = wallet == nil

	r.modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(15), Bottom: unit.Dp(15),
			Left: unit.Dp(15),
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
							return r.buttonClear.Layout(gtx, th)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(15)}.Layout),
					)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.Y = gtx.Dp(250)
					if wallet == nil {
						lbl := material.Label(th, unit.Sp(16), lang.Translate("Wallet is not opened."))
						return lbl.Layout(gtx)
					} else {
						if len(r.txItems) == 0 {
							lbl := material.Label(th, unit.Sp(16), lang.Translate("You don't have outgoing txs yet."))
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
	tx         wallet_manager.OutgoingTx
	buttonOpen *components.Button
}

func NewTxItem(tx wallet_manager.OutgoingTx) *TxItem {
	openIcon, _ := widget.NewIcon(icons.ActionOpenInBrowser)

	textColor := color.NRGBA{R: 0, G: 0, B: 0, A: 100}
	textHoverColor := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	buttonOpen := components.NewButton(components.ButtonStyle{
		Icon:           openIcon,
		TextColor:      textColor,
		HoverTextColor: &textHoverColor,
	})

	return &TxItem{
		tx:         tx,
		buttonOpen: buttonOpen,
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

	if item.buttonOpen.Clicked() {
		go open.Run(fmt.Sprintf("https://explorer.dero.io/tx/%s", txId))
	}

	r := op.Record(gtx.Ops)
	dims := layout.Inset{
		Top: unit.Dp(5), Bottom: unit.Dp(5),
		Left: unit.Dp(5), Right: unit.Dp(5),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Spacing:   layout.SpaceBetween,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), utils.ReduceTxId(txId))
						lbl.Font.Weight = font.Bold
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), status)
						return lbl.Layout(gtx)
					}),
				)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), fmt.Sprint(item.tx.BlockHeight.Int64))
						lbl.Alignment = text.End
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), lang.TimeAgo(date))
						lbl.Alignment = text.End
						return lbl.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return item.buttonOpen.Layout(gtx, th)
			}),
		)
	})
	c := r.Stop()

	paint.FillShape(gtx.Ops,
		color.NRGBA{A: 10},
		clip.RRect{
			Rect: image.Rectangle{Max: dims.Size},
			SE:   gtx.Dp(5), SW: gtx.Dp(5),
			NW: gtx.Dp(5), NE: gtx.Dp(5),
		}.Op(gtx.Ops),
	)

	c.Add(gtx.Ops)
	return dims
}
