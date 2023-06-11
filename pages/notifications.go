package pages

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/ui/components"
)

var SuccessModalInstance *components.NotificationModal
var ErrorModalInstance *components.NotificationModal
var InfoModalInstance *components.NotificationModal

func LoadNotificationsInstance() {
	w := app_instance.Window

	ErrorModalInstance = components.NewNotificationErrorModal(w)
	SuccessModalInstance = components.NewNotificationSuccessModal(w)
	InfoModalInstance := components.NewNotificationInfoModal(w)

	router := app_instance.Router
	router.PushLayout(func(gtx layout.Context, th *material.Theme) {
		SuccessModalInstance.Layout(gtx, th)
		ErrorModalInstance.Layout(gtx, th)
		InfoModalInstance.Layout(gtx, th)
	})
}
