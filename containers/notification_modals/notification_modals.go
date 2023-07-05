package notification_modals

import (
	"time"

	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/components"
)

var SuccessInstance *components.NotificationModal
var ErrorInstance *components.NotificationModal
var InfoInstance *components.NotificationModal
var CLOSE_AFTER_DEFAULT = 3 * time.Second

func LoadInstance() {
	ErrorInstance = components.NewNotificationErrorModal()
	SuccessInstance = components.NewNotificationSuccessModal()
	InfoInstance = components.NewNotificationInfoModal()

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 100,
		Layout: func(gtx layout.Context, th *material.Theme) {
			SuccessInstance.Layout(gtx, th)
			ErrorInstance.Layout(gtx, th)
			InfoInstance.Layout(gtx, th)
		},
	})
}
