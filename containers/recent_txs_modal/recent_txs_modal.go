package recent_txs_modal

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/components"
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
			openedWallet := wallet_manager.Instance.OpenedWallet
			if openedWallet == nil {
				lbl := material.Label(th, unit.Sp(14), "Wallet is not connected.")
				return lbl.Layout(gtx)
			} else {
				return r.listStyle.Layout(gtx, 0, func(gtx layout.Context, index int) layout.Dimensions {
					return layout.Dimensions{}
				})
			}
		})
	})
}
