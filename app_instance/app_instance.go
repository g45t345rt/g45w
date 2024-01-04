package app_instance

import (
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/unit"
	"gioui.org/x/explorer"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
)

var Window *app.Window
var Router *router.Router
var Explorer *explorer.Explorer

func Load() {
	minSize := image.Pt(375, 480)
	maxSize := image.Pt(480, 800)
	size := image.Pt(375, 625)

	Window = app.NewWindow(
		app.Title(settings.Name),
		app.MinSize(unit.Dp(minSize.X), unit.Dp(minSize.Y)),
		app.Size(unit.Dp(size.X), unit.Dp(size.Y)),
		app.MaxSize(unit.Dp(maxSize.X), unit.Dp(maxSize.Y)),
		app.PortraitOrientation.Option(),
		app.NavigationColor(color.NRGBA{A: 0}),
	)

	Explorer = explorer.NewExplorer(Window)
	Router = router.NewRouter()
}
