package containers

import (
	"github.com/g45t345rt/g45w/containers/bottom_bar"
	"github.com/g45t345rt/g45w/containers/build_tx_modal"
	"github.com/g45t345rt/g45w/containers/confirm_modal"
	"github.com/g45t345rt/g45w/containers/image_modal"
	"github.com/g45t345rt/g45w/containers/listselect_modal"
	"github.com/g45t345rt/g45w/containers/node_status_bar"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/containers/password_modal"
	"github.com/g45t345rt/g45w/containers/prompt_modal"
	"github.com/g45t345rt/g45w/containers/qrcode_scan_modal"
	"github.com/g45t345rt/g45w/containers/recent_txs_modal"
)

func Load() {
	bottom_bar.LoadInstance()
	node_status_bar.LoadInstance()
	notification_modals.LoadInstance()
	recent_txs_modal.LoadInstance()
	build_tx_modal.LoadInstance()
	image_modal.LoadInstance()
	qrcode_scan_modal.LoadInstance()
	confirm_modal.LoadInstance()
	password_modal.LoadInstance()
	prompt_modal.LoadInstance()
	listselect_modal.LoadInstance()
}
