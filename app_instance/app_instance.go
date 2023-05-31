package app_instance

import (
	"gioui.org/app"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/pages"
	"github.com/g45t345rt/g45w/router"
)

type AppInstance struct {
	Window    *app.Window
	Theme     *material.Theme
	Router    *router.Router
	BottomBar *pages.BottomBar
}

var Current *AppInstance
