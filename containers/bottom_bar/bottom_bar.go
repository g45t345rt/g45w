package bottom_bar

import (
	"fmt"
	"image"
	"image/color"

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
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/g45t345rt/g45w/wallet_manager"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type BottomBar struct {
	buttonWallet   *BottomBarButton
	buttonTxs      *BottomBarButton
	buttonSettings *BottomBarButton
	buttonClose    *BottomBarButton
	buttonNode     *BottomBarButton
	router         *router.Router

	confirmClose *components.Confirm
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
	w := app_instance.Window
	router := app_instance.Router
	th := app_instance.Theme

	textColor := color.NRGBA{R: 0, G: 0, B: 0, A: 100}
	textHoverColor := color.NRGBA{R: 0, G: 0, B: 0, A: 255} //f32color.Hovered(textColor)

	animScale := float32(.95)
	walletIcon, _ := widget.NewIcon(icons.ActionAccountBalanceWallet)
	buttonWallet := components.NewButton(components.ButtonStyle{
		Icon:           walletIcon,
		TextColor:      textColor,
		HoverTextColor: &textHoverColor,
		Animation:      components.NewButtonAnimationScale(animScale),
	})

	txsIcon, _ := widget.NewIcon(icons.ActionHistory)
	buttonTxs := components.NewButton(components.ButtonStyle{
		Icon:           txsIcon,
		TextColor:      textColor,
		HoverTextColor: &textHoverColor,
		Animation:      components.NewButtonAnimationScale(animScale),
	})

	settingsIcon, _ := widget.NewIcon(icons.ActionSettingsApplications)
	buttonSettings := components.NewButton(components.ButtonStyle{
		Icon:           settingsIcon,
		TextColor:      textColor,
		HoverTextColor: &textHoverColor,
		Animation:      components.NewButtonAnimationScale(animScale),
	})

	closeIcon, _ := widget.NewIcon(icons.ActionExitToApp)
	buttonClose := components.NewButton(components.ButtonStyle{
		Icon:           closeIcon,
		TextColor:      textColor,
		HoverTextColor: &textHoverColor,
		Animation:      components.NewButtonAnimationScale(animScale),
	})

	nodeIcon, _ := widget.NewIcon(icons.ActionDNS)
	buttonNode := components.NewButton(components.ButtonStyle{
		Icon:           nodeIcon,
		TextColor:      textColor,
		HoverTextColor: &textHoverColor,
		Animation:      components.NewButtonAnimationScale(animScale),
	})

	confirmClose := components.NewConfirm(w, "Closing current wallet?", th, layout.Center)
	router.PushLayout(func(gtx layout.Context, th *material.Theme) {
		confirmClose.Layout(gtx, th)
	})

	bottomBar := &BottomBar{
		buttonWallet:   NewBottomBarButton(buttonWallet),
		buttonTxs:      NewBottomBarButton(buttonTxs),
		buttonSettings: NewBottomBarButton(buttonSettings),
		buttonClose:    NewBottomBarButton(buttonClose),
		buttonNode:     NewBottomBarButton(buttonNode),
		confirmClose:   confirmClose,
		router:         router,
	}
	Instance = bottomBar
	return bottomBar
}

func (b *BottomBar) SetButtonActive(tag interface{}) {
	b.buttonSettings.button.Focused = false
	b.buttonClose.button.Focused = false
	b.buttonWallet.button.Focused = false
	b.buttonNode.button.Focused = false
	b.buttonTxs.button.Focused = false

	switch tag {
	case BUTTON_SETTINGS:
		b.buttonSettings.button.Focused = true
	case BUTTON_CLOSE:
		b.buttonClose.button.Focused = true
	case BUTTON_WALLET:
		b.buttonWallet.button.Focused = true
	case BUTTON_NODE:
		b.buttonNode.button.Focused = true
	case BUTTON_TXS:
		b.buttonTxs.button.Focused = true
	}
}

func (b *BottomBar) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.Rect{
		Max: gtx.Constraints.Max,
	}.Op())

	BottomBarTopWallet{}.Layout(gtx, th)

	if wallet_manager.Instance.OpenedWallet != nil {
		b.buttonClose.button.Disabled = false
		if b.buttonClose.button.Clickable.Clicked() {
			b.confirmClose.SetVisible(true)
		}
	} else {
		b.buttonClose.button.Disabled = true
	}

	if b.confirmClose.ClickedYes() {
		b.router.SetCurrent(app_instance.PAGE_WALLET_SELECT)
		wallet_manager.Instance.OpenedWallet = nil
	}

	if b.buttonNode.button.Clickable.Clicked() {
		b.router.SetCurrent(app_instance.PAGE_NODE)
	}

	if b.buttonWallet.button.Clickable.Clicked() {
		if b.router.Current == app_instance.PAGE_WALLET {
			b.router.SetCurrent(app_instance.PAGE_WALLET_SELECT)
		} else {
			b.router.SetCurrent(app_instance.PAGE_WALLET)
		}
	}

	if b.buttonSettings.button.Clickable.Clicked() {
		b.router.SetCurrent(app_instance.PAGE_SETTINGS)
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
				return b.buttonSettings.Layout(gtx, th)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return b.buttonTxs.Layout(gtx, th)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return b.buttonWallet.Layout(gtx, th)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return b.buttonNode.Layout(gtx, th)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return b.buttonClose.Layout(gtx, th)
			}),
		)
	})
}

type BottomBarButton struct {
	button *components.Button
}

func NewBottomBarButton(button *components.Button) *BottomBarButton {
	return &BottomBarButton{
		button: button,
	}
}

func (b *BottomBarButton) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	iconSize := unit.Dp(45)
	gtx.Constraints.Min.X = gtx.Dp(iconSize)
	gtx.Constraints.Min.Y = gtx.Dp(iconSize)

	if b.button.Focused {
		b.button.Style.TextColor = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	} else {
		// important scale down instead of up to avoid blurry icon
		scale := f32.Affine2D{}.Scale(f32.Pt(float32(iconSize)/2, float32(iconSize)/2), f32.Pt(.7, .7))
		defer op.Affine(scale).Push(gtx.Ops).Pop()

		b.button.Style.TextColor = color.NRGBA{R: 0, G: 0, B: 0, A: 100}
	}

	return b.button.Layout(gtx, th)
}

type BottomBarTopWallet struct{}

func (b BottomBarTopWallet) Layout(gtx layout.Context, th *material.Theme) {
	openedWallet := wallet_manager.Instance.OpenedWallet
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
			lbl := material.Label(th, unit.Sp(14), fmt.Sprintf("Wallet [%s]", walletName))
			lbl.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			lbl.Font.Weight = font.Bold
			return lbl.Layout(gtx)
		})
		c := r.Stop()

		x := float32(dims.Size.X / 2)
		y := float32(dims.Size.Y / 2)
		offset := f32.Pt(-x, -y)
		defer op.Affine(f32.Affine2D{}.Offset(offset)).Push(gtx.Ops).Pop()

		paint.FillShape(gtx.Ops, color.NRGBA{A: 255}, clip.UniformRRect(
			image.Rect(0, 0, dims.Size.X, dims.Size.Y),
			gtx.Dp(5),
		).Op(gtx.Ops))

		c.Add(gtx.Ops)
		return layout.Dimensions{}
	})
}
