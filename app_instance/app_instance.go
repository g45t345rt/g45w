package app_instance

import (
	"gioui.org/app"
	"gioui.org/x/explorer"
	"github.com/g45t345rt/g45w/router"
)

var Window *app.Window
var Router *router.Router
var Explorer *explorer.Explorer

const (
	PAGE_SETTINGS      = "page_settings"
	PAGE_NODE          = "page_node"
	PAGE_WALLET        = "page_wallet"
	PAGE_WALLET_SELECT = "page_wallet_select"
)
