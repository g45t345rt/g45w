package main

import (
	"image/color"
	"log"
	"os"
	"time"

	"eliasnaur.com/font/roboto/robotobold"
	"eliasnaur.com/font/roboto/robotoregular"
	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/node"
	"github.com/g45t345rt/g45w/pages"
	page_node "github.com/g45t345rt/g45w/pages/node"
	page_settings "github.com/g45t345rt/g45w/pages/settings"
	page_wallet "github.com/g45t345rt/g45w/pages/wallet"
	page_wallet_select "github.com/g45t345rt/g45w/pages/wallet_select"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
)

func main() {
	err := settings.NewSettings().LoadSettings()
	if err != nil {
		log.Fatal(err)
	}

	err = wallet_manager.NewWalletManager().LoadWallets()
	if err != nil {
		log.Fatal(err)
	}

	err = node.NewNode().Start()
	if err != nil {
		log.Fatal(err)
	}

	// window
	minSizeX := unit.Dp(375)
	minSizeY := unit.Dp(600)
	maxSizeX := unit.Dp(500)
	maxSizeY := unit.Dp(1000)

	window := app.NewWindow(
		app.Title("G45W"),
		app.MinSize(minSizeX, minSizeY),
		app.Size(minSizeX, minSizeY),
		app.MaxSize(maxSizeX, maxSizeY),
		app.StatusColor(color.NRGBA{R: 255, A: 255}),
		app.PortraitOrientation.Option(),
		app.NavigationColor(color.NRGBA{R: 255, G: 255, A: 255}),
	)

	// font
	fontCollection, err := loadFontCollection()
	if err != nil {
		log.Fatal(err)
	}

	// theme
	theme := material.NewTheme(fontCollection)
	theme.WithPalette(material.Palette{
		Fg:         utils.HexColor(0x000000),
		Bg:         utils.HexColor(0xffffff),
		ContrastBg: utils.HexColor(0x3f51b5),
		ContrastFg: utils.HexColor(0xffffff),
	})

	// main router
	router := router.NewRouter()

	// app instance to give guick access to every package
	app_instance.Current = &app_instance.AppInstance{
		Window:    window,
		Theme:     theme,
		Router:    router,
		BottomBar: pages.NewBottomBar(router, theme),
	}

	router.Add("page_settings", page_settings.NewPage())
	router.Add("page_node", page_node.NewPage())
	router.Add("page_wallet", page_wallet.NewPage())
	router.Add("page_wallet_select", page_wallet_select.NewPage())
	router.SetPrimary("page_wallet_select")

	go func() {
		err := runApp(window, router, theme)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func loadFontCollection() ([]font.FontFace, error) {
	robotoRegular, err := opentype.Parse(robotoregular.TTF)
	if err != nil {
		return nil, err
	}

	robotoBold, err := opentype.Parse(robotobold.TTF)
	if err != nil {
		return nil, err
	}

	fontCollection := []font.FontFace{}
	//gofont.Collection()
	fontCollection = append(fontCollection, font.FontFace{Font: font.Font{}, Face: robotoRegular})
	fontCollection = append(fontCollection, font.FontFace{Font: font.Font{Weight: font.Bold}, Face: robotoBold})
	return fontCollection, nil
}

func runApp(window *app.Window, router *router.Router, th *material.Theme) error {
	var ops op.Ops

	// 1s ticker to update node status and topbar...
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case e := <-window.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				router.Layout(gtx, th)
				e.Frame(gtx.Ops)
			}
		case <-ticker.C:
			window.Invalidate()
		}
	}
}
