package main

import (
	"image/color"
	"log"
	"os"

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
	expl "gioui.org/x/explorer"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/app_data"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/containers/bottom_bar"
	"github.com/g45t345rt/g45w/containers/build_tx_modal.go"
	"github.com/g45t345rt/g45w/containers/node_status_bar"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/containers/recent_txs_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/node_manager"
	page_node "github.com/g45t345rt/g45w/pages/node"
	page_settings "github.com/g45t345rt/g45w/pages/settings"
	page_wallet "github.com/g45t345rt/g45w/pages/wallet"
	page_wallet_select "github.com/g45t345rt/g45w/pages/wallet_select"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
)

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

func runApp(th *material.Theme) error {
	var ops op.Ops

	window := app_instance.Window
	router := app_instance.Router
	explorer := app_instance.Explorer

	loadState := NewLoadState(window)

	go func() {
		loadState.logoSplash.animation.Start()
		loadState.SetStatus("Initiating", nil) // don't use lang.Translate - lang is not loaded

		err := lang.Load()
		if err != nil {
			loadState.SetStatus("", err)
			return
		}

		loadState.SetStatus(lang.Translate("Loading settings"), nil)
		err = settings.Load()
		if err != nil {
			loadState.SetStatus("", err)
			return
		}

		loadState.SetStatus(lang.Translate("Loading lookup table"), nil)
		walletapi.Initialize_LookupTable(1, 1<<21)

		loadState.SetStatus(lang.Translate("Loading app data"), nil)
		err = app_data.Load()
		if err != nil {
			loadState.SetStatus("", err)
			return
		}

		loadState.SetStatus(lang.Translate("Loading wallets"), nil)
		err = wallet_manager.Load()
		if err != nil {
			loadState.SetStatus("", err)
			return
		}

		loadState.SetStatus(lang.Translate("Checking node"), nil)
		err = node_manager.Load()
		if err != nil {
			loadState.SetStatus("", err)
			return
		}

		loadState.SetStatus(lang.Translate("Loading pages"), nil)
		bottom_bar.LoadInstance()
		node_status_bar.LoadInstance()
		notification_modals.LoadInstance()
		recent_txs_modal.LoadInstance()
		build_tx_modal.LoadInstance()

		router.Add(app_instance.PAGE_NODE, page_node.New())
		router.Add(app_instance.PAGE_WALLET, page_wallet.New())
		router.Add(app_instance.PAGE_WALLET_SELECT, page_wallet_select.New())
		router.Add(app_instance.PAGE_SETTINGS, page_settings.New())
		router.SetCurrent(app_instance.PAGE_WALLET_SELECT)

		loadState.logoSplash.animation.Pause()
		loadState.SetStatus(lang.Translate("Done"), nil)
		loadState.Complete()
	}()

	for {
		e := <-window.Events()
		explorer.ListenEvents(e)
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			if loadState.loaded {
				router.Layout(gtx, th)
			} else {
				loadState.Layout(gtx, th)
			}

			e.Frame(gtx.Ops)
		}
	}
}

func main() {
	globals.Arguments["--testnet"] = false
	globals.Arguments["--debug"] = false
	globals.Arguments["--flog-level"] = nil
	globals.Arguments["--log-dir"] = nil
	globals.Arguments["--help"] = false
	globals.Arguments["--version"] = false

	// window
	minSizeX := unit.Dp(375)
	minSizeY := unit.Dp(600)
	maxSizeX := unit.Dp(500)
	maxSizeY := unit.Dp(1000)

	window := app.NewWindow(
		app.Title(settings.Name),
		app.MinSize(minSizeX, minSizeY),
		app.Size(minSizeX, minSizeY),
		app.MaxSize(maxSizeX, maxSizeY),
		app.PortraitOrientation.Option(),
		app.NavigationColor(color.NRGBA{A: 0}),
	)

	explorer := expl.NewExplorer(window)

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
	theme.FingerSize = 48

	// main router
	appRouter := router.NewRouter()

	// app instance to give guick access to every package
	app_instance.Window = window
	app_instance.Router = appRouter
	app_instance.Explorer = explorer

	go func() {
		err := runApp(theme)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}
