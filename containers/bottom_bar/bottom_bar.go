package bottom_bar

import (
	"image"
	"strings"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/confirm_modal"
	"github.com/g45t345rt/g45w/containers/recent_txs_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/pages"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/wallet_manager"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type BottomBar struct {
	ButtonWallet   *BottomBarButton
	ButtonTxs      *BottomBarButton
	ButtonSettings *BottomBarButton
	ButtonClose    *BottomBarButton
	ButtonNode     *BottomBarButton

	appRouter *router.Router
}

var Instance *BottomBar

const (
	BUTTON_WALLET   = "wallet"
	BUTTON_NODE     = "node"
	BUTTON_TXS      = "txs"
	BUTTON_CLOSE    = "close"
	BUTTON_SETTINGS = "settings"
)

func LoadInstance() *BottomBar {
	appRouter := app_instance.Router

	animScale := float32(.95)
	walletIcon, _ := widget.NewIcon(icons.ActionAccountBalanceWallet)
	buttonWallet := components.NewButton(components.ButtonStyle{
		Icon:      walletIcon,
		Animation: components.NewButtonAnimationScale(animScale),
	})

	txsIcon, _ := widget.NewIcon(icons.ActionHistory)
	buttonTxs := components.NewButton(components.ButtonStyle{
		Icon:      txsIcon,
		Animation: components.NewButtonAnimationScale(animScale),
	})

	settingsIcon, _ := widget.NewIcon(icons.ActionSettingsApplications)
	buttonSettings := components.NewButton(components.ButtonStyle{
		Icon:      settingsIcon,
		Animation: components.NewButtonAnimationScale(animScale),
	})

	closeIcon, _ := widget.NewIcon(icons.ActionExitToApp)
	buttonClose := components.NewButton(components.ButtonStyle{
		Icon:      closeIcon,
		Animation: components.NewButtonAnimationScale(animScale),
	})

	nodeIcon, _ := widget.NewIcon(icons.ActionDNS)
	buttonNode := components.NewButton(components.ButtonStyle{
		Icon:      nodeIcon,
		Animation: components.NewButtonAnimationScale(animScale),
	})

	bottomBar := &BottomBar{
		ButtonWallet:   NewBottomBarButton(buttonWallet),
		ButtonTxs:      NewBottomBarButton(buttonTxs),
		ButtonSettings: NewBottomBarButton(buttonSettings),
		ButtonClose:    NewBottomBarButton(buttonClose),
		ButtonNode:     NewBottomBarButton(buttonNode),
		appRouter:      appRouter,
	}
	Instance = bottomBar
	return bottomBar
}

func (b *BottomBar) SetButtonActive(tag interface{}) {
	b.ButtonSettings.Button.Focused = false
	b.ButtonClose.Button.Focused = false
	b.ButtonWallet.Button.Focused = false
	b.ButtonNode.Button.Focused = false
	b.ButtonTxs.Button.Focused = false

	switch tag {
	case BUTTON_SETTINGS:
		b.ButtonSettings.Button.Focused = true
	case BUTTON_CLOSE:
		b.ButtonClose.Button.Focused = true
	case BUTTON_WALLET:
		b.ButtonWallet.Button.Focused = true
	case BUTTON_NODE:
		b.ButtonNode.Button.Focused = true
	case BUTTON_TXS:
		b.ButtonTxs.Button.Focused = true
	}
}

func (b *BottomBar) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	bgColor := theme.Current.BottomBarBgColor
	paint.FillShape(gtx.Ops, bgColor, clip.Rect{
		Max: gtx.Constraints.Max,
	}.Op())

	BottomBarTopWallet{}.Layout(gtx, th)

	if wallet_manager.OpenedWallet != nil {
		b.ButtonClose.Button.Disabled = false
		if b.ButtonClose.Button.Clicked() {
			go func() {
				yesChan := confirm_modal.Instance.Open(confirm_modal.ConfirmText{
					Prompt: lang.Translate("Closing current wallet?"),
				})

				for yes := range yesChan {
					if yes {
						b.appRouter.SetCurrent(pages.PAGE_WALLET_SELECT)
						wallet_manager.CloseOpenedWallet()
					}
				}
			}()
		}
	} else {
		b.ButtonClose.Button.Disabled = true
	}

	if b.ButtonNode.Button.Clicked() {
		b.appRouter.SetCurrent(pages.PAGE_NODE)
	}

	if b.ButtonWallet.Button.Clicked() {
		if b.appRouter.Current == pages.PAGE_WALLET {
			b.appRouter.SetCurrent(pages.PAGE_WALLET_SELECT)
		} else {
			b.appRouter.SetCurrent(pages.PAGE_WALLET)
		}
	}

	if b.ButtonSettings.Button.Clicked() {
		b.appRouter.SetCurrent(pages.PAGE_SETTINGS)
	}

	if b.ButtonTxs.Button.Clicked() {
		recent_txs_modal.Instance.SetVisible(true)
	}

	return layout.Inset{
		Top: unit.Dp(20), Bottom: unit.Dp(20),
		Left: unit.Dp(30), Right: unit.Dp(30),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Spacing:   layout.SpaceBetween,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return b.ButtonSettings.Layout(gtx, th)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return b.ButtonTxs.Layout(gtx, th)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return b.ButtonWallet.Layout(gtx, th)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return b.ButtonNode.Layout(gtx, th)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return b.ButtonClose.Layout(gtx, th)
			}),
		)
	})
}

type BottomBarButton struct {
	Button *components.Button
}

func NewBottomBarButton(button *components.Button) *BottomBarButton {
	return &BottomBarButton{
		Button: button,
	}
}

func (b *BottomBarButton) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	iconSize := unit.Dp(45)
	gtx.Constraints.Min.X = gtx.Dp(iconSize)
	gtx.Constraints.Min.Y = gtx.Dp(iconSize)

	b.Button.Style.Colors = theme.Current.BottomButtonColors
	if b.Button.Focused {
		b.Button.Style.Colors.TextColor = theme.Current.BottomButtonSelectedColor
	} else {
		// important scale down instead of up to avoid blurry icon
		scale := f32.Affine2D{}.Scale(f32.Pt(float32(iconSize)/2, float32(iconSize)/2), f32.Pt(.7, .7))
		defer op.Affine(scale).Push(gtx.Ops).Pop()
	}

	return b.Button.Layout(gtx, th)
}

type BottomBarTopWallet struct{}

func (b BottomBarTopWallet) Layout(gtx layout.Context, th *material.Theme) {
	openedWallet := wallet_manager.OpenedWallet
	if openedWallet == nil {
		return
	}

	layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		r := op.Record(gtx.Ops)
		dims := layout.Inset{
			Top: unit.Dp(6), Bottom: unit.Dp(6),
			Left: unit.Dp(10), Right: unit.Dp(10),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			walletName := openedWallet.Info.Name
			text := lang.Translate("Wallet [{}]")
			text = strings.Replace(text, "{}", walletName, -1)
			lbl := material.Label(th, unit.Sp(14), text)
			lbl.Color = theme.Current.BottomBarWalletTextColor
			lbl.Font.Weight = font.Bold
			return lbl.Layout(gtx)
		})
		c := r.Stop()

		x := float32(dims.Size.X / 2)
		y := float32(dims.Size.Y / 2)
		offset := f32.Pt(-x, -y)
		defer op.Affine(f32.Affine2D{}.Offset(offset)).Push(gtx.Ops).Pop()

		bgColor := theme.Current.BottomBarWalletBgColor
		paint.FillShape(gtx.Ops, bgColor, clip.UniformRRect(
			image.Rect(0, 0, dims.Size.X, dims.Size.Y),
			gtx.Dp(5),
		).Op(gtx.Ops))

		c.Add(gtx.Ops)
		return layout.Dimensions{}
	})
}
