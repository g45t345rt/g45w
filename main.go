package main

import (
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
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/app_data"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/containers"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/node_manager"
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
	fontCollection = append(fontCollection, font.FontFace{Font: font.Font{}, Face: robotoRegular})
	fontCollection = append(fontCollection, font.FontFace{Font: font.Font{Weight: font.Bold}, Face: robotoBold})
	return fontCollection, nil
}

func loadPages(router *router.Router) {
	pageNode := page_node.New()
	pageWallet := page_wallet.New()
	pageWalletSelect := page_wallet_select.New()
	pageSettings := page_settings.New()

	router.Add(pages.PAGE_NODE, pageNode)
	router.Add(pages.PAGE_WALLET, pageWallet)
	router.Add(pages.PAGE_WALLET_SELECT, pageWalletSelect)
	router.Add(pages.PAGE_SETTINGS, pageSettings)
	router.SetCurrent(pages.PAGE_WALLET_SELECT)
}

func runApp() error {
	globals.Arguments["--testnet"] = false
	globals.Arguments["--debug"] = false
	globals.Arguments["--flog-level"] = nil
	globals.Arguments["--log-dir"] = nil
	globals.Arguments["--help"] = false
	globals.Arguments["--version"] = false

	var ops op.Ops
	app_instance.Load()
	window := app_instance.Window
	explorer := app_instance.Explorer
	router := app_instance.Router

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

		loadState.SetStatus(lang.Translate("Loading nodes"), nil)
		node_manager.Load() // don't check for error (e.g if current node connected successfully) and continue loading the app

		loadState.SetStatus(lang.Translate("Loading pages"), nil)
		containers.Load()
		loadPages(router)

		loadState.logoSplash.animation.Pause()
		loadState.SetStatus(lang.Translate("Done"), nil)
		loadState.Complete()
	}()

	fontCollection, err := loadFontCollection()
	if err != nil {
		log.Fatal(err)
	}

	th := material.NewTheme(fontCollection)
	th.WithPalette(material.Palette{
		Fg:         utils.HexColor(0x000000),
		Bg:         utils.HexColor(0xffffff),
		ContrastBg: utils.HexColor(0x3f51b5),
		ContrastFg: utils.HexColor(0xffffff),
	})
	th.FingerSize = 48

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
	go func() {
		err := runApp()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}
