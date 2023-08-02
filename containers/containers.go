package containers

import (
	"github.com/g45t345rt/g45w/containers/bottom_bar"
	"github.com/g45t345rt/g45w/containers/build_tx_modal.go"
	"github.com/g45t345rt/g45w/containers/node_status_bar"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/containers/recent_txs_modal"
)

func Load() {
	bottom_bar.LoadInstance()
	node_status_bar.LoadInstance()
	notification_modals.LoadInstance()
	recent_txs_modal.LoadInstance()
	build_tx_modal.LoadInstance()
}
