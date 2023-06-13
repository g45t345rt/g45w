package notification_modals

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/ui/components"
)

var SuccessInstance *components.NotificationModal
var ErrorInstance *components.NotificationModal
var InfoInstance *components.NotificationModal

func LoadInstance() {
	w := app_instance.Window

	ErrorInstance = components.NewNotificationErrorModal(w)
	SuccessInstance = components.NewNotificationSuccessModal(w)
	InfoInstance = components.NewNotificationInfoModal(w)

	app_instance.Router.PushLayout(func(gtx layout.Context, th *material.Theme) {
		SuccessInstance.Layout(gtx, th)
		ErrorInstance.Layout(gtx, th)
		InfoInstance.Layout(gtx, th)
	})
}
