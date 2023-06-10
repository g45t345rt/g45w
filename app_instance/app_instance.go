package app_instance

import (
	"gioui.org/app"
	"gioui.org/widget/material"
	"gioui.org/x/explorer"
	"github.com/g45t345rt/g45w/router"
)

type AppInstance struct {
	Window   *app.Window
	Theme    *material.Theme
	Router   *router.Router
	Explorer *explorer.Explorer
}

var Window *app.Window
var Theme *material.Theme
var Router *router.Router
var Explorer *explorer.Explorer

var Current *AppInstance
