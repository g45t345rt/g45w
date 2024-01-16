package notification_modal

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

var Instance *components.NotificationModal
var CLOSE_AFTER_DEFAULT = 3 * time.Second

type Type string

var (
	ERROR   Type = "error"
	SUCCESS Type = "success"
	INFO    Type = "info"
)

var currentType Type
var iconError *widget.Icon
var iconSuccess *widget.Icon
var iconInfo *widget.Icon

type Params struct {
	Type       Type
	Title      string
	Text       string
	CloseAfter time.Duration
}

func LoadInstance() {
	window := app_instance.Window

	iconError, _ = widget.NewIcon(icons.AlertError)
	iconSuccess, _ = widget.NewIcon(icons.ActionCheckCircle)
	iconInfo, _ = widget.NewIcon(icons.ActionInfo)

	Instance = components.NewNotificationModal(window, components.NotificationStyle{
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

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 100,
		Layout: func(gtx layout.Context, th *material.Theme) {
			switch currentType {
			case SUCCESS:
				Instance.Style.Colors = theme.Current.NotificationSuccessColors
			case ERROR:
				Instance.Style.Colors = theme.Current.NotificationErrorColors
			case INFO:
				Instance.Style.Colors = theme.Current.NotificationInfoColors
			}

			Instance.Layout(gtx, th)
		},
	})
}

func Open(params Params) {
	Instance.SetVisible(false, 0)

	currentType = params.Type
	switch params.Type {
	case SUCCESS:
		Instance.Style.Icon = iconSuccess
	case ERROR:
		Instance.Style.Icon = iconError
	case INFO:
		Instance.Style.Icon = iconInfo
	}

	Instance.SetText(params.Title, params.Text)
	Instance.SetVisible(true, params.CloseAfter)
}
