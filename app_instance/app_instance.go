package app_instance

import (
	"gioui.org/app"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/blockchain"
	"github.com/g45t345rt/g45w/pages"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
)

type AppInstance struct {
	Window    *app.Window
	Theme     *material.Theme
	Router    *router.Router
	BottomBar *pages.BottomBar
	Chain     *blockchain.Blockchain
	Settings  *settings.Settings
}

var Current *AppInstance
