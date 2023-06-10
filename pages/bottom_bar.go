package pages

import (
	"image/color"

	"gioui.org/f32"
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

var BottomBarInstance *BottomBar

func LoadBottomBarInstance() *BottomBar {
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
	BottomBarInstance = bottomBar
	return bottomBar
}

func (b *BottomBar) SetButtonActive(tag string) {
	b.buttonSettings.button.Focused = false
	b.buttonClose.button.Focused = false
	b.buttonWallet.button.Focused = false
	b.buttonNode.button.Focused = false
	b.buttonTxs.button.Focused = false

	switch tag {
	case "settings":
		b.buttonSettings.button.Focused = true
	case "close":
		b.buttonClose.button.Focused = true
	case "wallet":
		b.buttonWallet.button.Focused = true
	case "node":
		b.buttonNode.button.Focused = true
	case "txs":
		b.buttonTxs.button.Focused = true
	}
}

func (b *BottomBar) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.Rect{
		Max: gtx.Constraints.Max,
	}.Op())

	if wallet_manager.Instance.OpenedWallet != nil {
		b.buttonClose.button.Disabled = false
		if b.buttonClose.button.Clickable.Clicked() {
			b.confirmClose.SetVisible(true)
		}
	} else {
		b.buttonClose.button.Disabled = true
	}

	if b.confirmClose.ClickedYes() {
		b.router.SetCurrent("page_wallet_select")
		wallet_manager.Instance.OpenedWallet = nil
	}

	if b.buttonNode.button.Clickable.Clicked() {
		b.router.SetCurrent("page_node")
	}

	if b.buttonWallet.button.Clickable.Clicked() {
		if b.router.Current == "page_wallet" {
			b.router.SetCurrent("page_wallet_select")
		} else {
			b.router.SetCurrent("page_wallet")
		}
	}

	if b.buttonSettings.button.Clickable.Clicked() {
		b.router.SetCurrent("page_settings")
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
