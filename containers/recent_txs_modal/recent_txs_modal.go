package recent_txs_modal

import (
	"fmt"
	"image/color"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
)

type RecentTxsModal struct {
	modal *components.Modal
	list  *widget.List
}

var Instance *RecentTxsModal

func LoadInstance() {
	w := app_instance.Window
	modal := components.NewModal(w, components.ModalStyle{
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

	list := new(widget.List)
	list.Axis = layout.Vertical

	Instance = &RecentTxsModal{
		modal: modal,
		list:  list,
	}

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 2,
		Layout:    Instance.layout,
	})
}

func (r *RecentTxsModal) SetVisible(visible bool) {
	r.modal.SetVisible(visible)
}

func (r *RecentTxsModal) layout(gtx layout.Context, th *material.Theme) {
	var txItems []TxItem
	for i := 0; i < 10; i++ {
		txItems = append(txItems, TxItem{
			TxId:          "2fb45948c17337446ac54ec7644df5335d7a9b55454f4d286a210af937e34bbe",
			Confirmations: 23,
			Date:          time.Now(),
			Status:        "Checking transaction...",
		})
	}

	r.modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(15), Bottom: unit.Dp(15),
			Left: unit.Dp(15),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(20), fmt.Sprintf("%s (%d)", lang.Translate("Recent Transactions"), len(txItems)))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.Y = gtx.Dp(250)
					openedWallet := wallet_manager.OpenedWallet
					if openedWallet == nil {
						lbl := material.Label(th, unit.Sp(16), lang.Translate("Wallet is not opened."))
						return lbl.Layout(gtx)
					} else {
						listStyle := material.List(th, r.list)
						listStyle.AnchorStrategy = material.Overlay

						return listStyle.Layout(gtx, len(txItems), func(gtx layout.Context, index int) layout.Dimensions {
							bottom := 0
							if index < len(txItems)-1 {
								bottom = 10
							}

							return layout.Inset{
								Bottom: unit.Dp(bottom), Right: unit.Dp(15),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return txItems[index].Layout(gtx, th)
							})
						})
					}
				}),
			)
		})
	})
}

type TxItem struct {
	TxId          string
	Confirmations uint64
	Date          time.Time
	Status        string
}

func (tx *TxItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					txId := utils.ReduceTxId(tx.TxId)
					lbl := material.Label(th, unit.Sp(16), txId)
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), tx.Status)
					return lbl.Layout(gtx)
				}),
			)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					confirmations := fmt.Sprintf("%d confirmations", tx.Confirmations)
					lbl := material.Label(th, unit.Sp(16), confirmations)
					lbl.Alignment = text.End
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), tx.Date.Format("2006-01-02"))
					lbl.Alignment = text.End
					return lbl.Layout(gtx)
				}),
			)
		}),
	)
}
