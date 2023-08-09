package notification_modals

import (
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var SuccessInstance *components.NotificationModal
var ErrorInstance *components.NotificationModal
var InfoInstance *components.NotificationModal
var CLOSE_AFTER_DEFAULT = 3 * time.Second

func LoadInstance() {
	window := app_instance.Window

	iconError, _ := widget.NewIcon(icons.AlertError)
	ErrorInstance = components.NewNotificationModal(window, components.NotificationStyle{
		Direction:  layout.N,
		OuterInset: layout.UniformInset(unit.Dp(10)),
		InnerInset: layout.Inset{
			Top: unit.Dp(10), Bottom: unit.Dp(10),
			Left: unit.Dp(15), Right: unit.Dp(15),
		},
		Rounded:   components.UniformRounded(unit.Dp(10)),
		Icon:      iconError,
		Animation: components.NewModalAnimationDown(),
	})

	iconSuccess, _ := widget.NewIcon(icons.ActionCheckCircle)
	SuccessInstance = components.NewNotificationModal(window, components.NotificationStyle{
		Direction:  layout.N,
		OuterInset: layout.UniformInset(unit.Dp(10)),
		InnerInset: layout.Inset{
			Top: unit.Dp(10), Bottom: unit.Dp(10),
			Left: unit.Dp(15), Right: unit.Dp(15),
		},
		Rounded:   components.UniformRounded(unit.Dp(10)),
		Icon:      iconSuccess,
		Animation: components.NewModalAnimationDown(),
	})

	iconInfo, _ := widget.NewIcon(icons.ActionInfo)
	InfoInstance = components.NewNotificationModal(window, components.NotificationStyle{
		Direction:  layout.N,
		OuterInset: layout.UniformInset(unit.Dp(10)),
		InnerInset: layout.Inset{
			Top: unit.Dp(10), Bottom: unit.Dp(10),
			Left: unit.Dp(15), Right: unit.Dp(15),
		},
		Rounded:   components.UniformRounded(unit.Dp(10)),
		Icon:      iconInfo,
		Animation: components.NewModalAnimationDown(),
	})

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 100,
		Layout: func(gtx layout.Context, th *material.Theme) {
			SuccessInstance.Style.Colors = theme.Current.NotificationSuccessColors
			SuccessInstance.Layout(gtx, th)
			ErrorInstance.Style.Colors = theme.Current.NotificationErrorColors
			ErrorInstance.Layout(gtx, th)
			InfoInstance.Style.Colors = theme.Current.NotificationInfoColors
			InfoInstance.Layout(gtx, th)
		},
	})
}
