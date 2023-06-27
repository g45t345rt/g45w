package recent_txs_modal

import (
	"fmt"
	"image/color"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
)

type RecentTxsModal struct {
	modal     *components.Modal
	listStyle material.ListStyle
}

var Instance *RecentTxsModal

func LoadInstance() {
	w := app_instance.Window
	th := app_instance.Theme
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
	listStyle := material.List(th, list)
	listStyle.AnchorStrategy = material.Overlay

	Instance = &RecentTxsModal{
		modal:     modal,
		listStyle: listStyle,
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
	r.modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(18), "Recent Transactions")
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					openedWallet := wallet_manager.OpenedWallet
					if openedWallet == nil {
						lbl := material.Label(th, unit.Sp(14), "Wallet is not opened.")
						return lbl.Layout(gtx)
					} else {
						txItems := []TxItem{
							{
								TxId:          "b0a555db9dcac8d7bb2d9ccc27ade33f81ed6e4a283c50e28d77eb208cfa7ff2",
								Confirmations: 5,
								Date:          time.Now(),
								Status:        "Checking transaction...",
							},
						}

						return r.listStyle.Layout(gtx, len(txItems), func(gtx layout.Context, index int) layout.Dimensions {
							return txItems[index].Layout(gtx, th)
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
					lbl := material.Label(th, unit.Sp(14), txId)
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), tx.Status)
					return lbl.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					confirmations := fmt.Sprintf("%d confirmations", tx.Confirmations)
					lbl := material.Label(th, unit.Sp(14), confirmations)
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), tx.Date.Format("2006-01-02"))
					return lbl.Layout(gtx)
				}),
			)
		}),
	)
}
