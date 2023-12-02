package main

import (
	// first we are going to import some boilerplate
	"log"
	"os"

	// then we are going to import gio ui
	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
	"gioui.org/x/camera"

	// add android permissions
	_ "gioui.org/app/permission/camera"
	_ "gioui.org/app/permission/networkstate"
	_ "gioui.org/app/permission/storage"

	// support webp image decode
	//_ "github.com/chai2010/webp"

	//then we are going to import the dero project globals
	"github.com/deroproject/derohe/globals"

	// then we are going to import g45's wallet repo
	"github.com/g45t345rt/g45w/app_db"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/containers"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/lookup_table"
	"github.com/g45t345rt/g45w/node_manager"
	"github.com/g45t345rt/g45w/pages"
	page_node "github.com/g45t345rt/g45w/pages/node"
	page_settings "github.com/g45t345rt/g45w/pages/settings"
	page_wallet "github.com/g45t345rt/g45w/pages/wallet"
	page_wallet_select "github.com/g45t345rt/g45w/pages/wallet_select"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/theme"
)

// now, obviously we are going to want a font set
func loadFontCollection() ([]font.FontFace, error) {
	// universal fonts from https://github.com/satbyy/go-noto-universal

	RobotoRegularTTF, err := assets.GetFont("Roboto-Regular.ttf")
	if err != nil {
		return nil, err
	}

	RobotoRegular, err := opentype.Parse(RobotoRegularTTF)
	if err != nil {
		return nil, err
	}

	RobotoBoldTTF, err := assets.GetFont("Roboto-Bold.ttf")
	if err != nil {
		return nil, err
	}

	RobotoBold, err := opentype.Parse(RobotoBoldTTF)
	if err != nil {
		return nil, err
	}

	fontCollection := []font.FontFace{}
	fontCollection = append(fontCollection, font.FontFace{Font: font.Font{}, Face: RobotoRegular})
	fontCollection = append(fontCollection, font.FontFace{Font: font.Font{Weight: font.Bold}, Face: RobotoBold})
	return fontCollection, nil
}

// now we are going to make all of our pages
func loadPages(router *router.Router) {
	pageNode := page_node.New()
	pageWallet := page_wallet.New()
	pageWalletSelect := page_wallet_select.New()
	pageSettings := page_settings.New()

	// then we are going to add pages to our router
	router.Add(pages.PAGE_NODE, pageNode)
	router.Add(pages.PAGE_WALLET, pageWallet)
	router.Add(pages.PAGE_WALLET_SELECT, pageWalletSelect)
	router.Add(pages.PAGE_SETTINGS, pageSettings)

	// and now we are goign to set the current page
	router.SetCurrent(pages.PAGE_WALLET_SELECT)
}

func runApp() error {
	//  fundamentally, we are are using derohe's api
	globals.Arguments["--testnet"] = false
	globals.Arguments["--debug"] = false
	globals.Arguments["--flog-level"] = nil
	globals.Arguments["--log-dir"] = nil
	globals.Arguments["--help"] = false
	globals.Arguments["--version"] = false

	// so let's turn on the DERO mainnet network
	globals.InitNetwork() // this func assign mainnet/testnet config depending on globals.Arguments["--testnet"] value

	var ops op.Ops

	// now let's make the window
	app_instance.Load()
	window := app_instance.Window
	explorer := app_instance.Explorer
	router := app_instance.Router

	// let's load our theme
	theme.LoadImages()

	// load a window with splash
	loadState := NewLoadState(window)

	// now let's make a sequence for loading settings into memory
	go func() {

		// start the splash screen
		loadState.logoSplash.animation.Start()

		// start app status
		loadState.SetStatus("Initiating", nil) // don't use lang.Translate - lang is not loaded

		// load languages
		err := lang.Load()
		if err != nil {
			loadState.SetStatus("", err)
			return
		}

		// soft test language translation
		// and set status
		loadState.SetStatus(lang.Translate("Loading settings"), nil)

		// load settings into memory
		// it might be easier to call them configs...
		// but load them anyway
		err = settings.Load()
		if err != nil {
			loadState.SetStatus("", err)
			return
		}

		// set status
		loadState.SetStatus(lang.Translate("Loading lookup table"), nil)
		//walletapi.Initialize_LookupTable(1, 1<<21)

		// and load the precompiled lookup table into memory
		err = lookup_table.Load()
		if err != nil {
			loadState.SetStatus("", err)
			return
		}

		// set status
		loadState.SetStatus(lang.Translate("Loading app data"), nil)

		// load sqlite database into memory
		err = app_db.Load()
		if err != nil {
			loadState.SetStatus("", err)
			return
		}

		// now load the node connection
		node_manager.Load() // don't check for error (e.g if current node connected successfully) and continue loading the app

		// set status
		loadState.SetStatus(lang.Translate("Loading pages"), nil)

		// load all of the divisions/containers of the applicaiton
		containers.Load()

		// load the router into the page loader
		loadPages(router)

		// update splash screen
		loadState.logoSplash.animation.Pause()

		// set status
		loadState.SetStatus(lang.Translate("Done"), nil)

		// mission complete: settings are loaded
		loadState.Complete()
	}()

	// now load the fonts
	fontCollection, err := loadFontCollection()
	if err != nil {
		log.Fatal(err)
	}

	// set the theme
	th := material.NewTheme()
	// tune the fonts
	th.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(fontCollection))
	// set fingerprint size
	th.FingerSize = 48

	// now this is the state/events engine
	for {
		// load all events into memory
		e := <-window.Events()

		// pass all events through the app explorer
		// window_explorer != blockchain_explorer
		explorer.ListenEvents(e)

		// pass all events through the camera
		camera.ListenEvents(e)

		// now let's make a switch for events
		switch e := e.(type) {
		// if we run into a big problem, run an error
		case system.DestroyEvent:
			return e.Err
			//if we have a frame change...
		case system.FrameEvent:
			// establish what the "new context" is going to be
			gtx := layout.NewContext(&ops, e)

			// paint the window as per theme
			paint.Fill(gtx.Ops, theme.Current.BgColor)
			th.Bg = theme.Current.BgColor
			th.Fg = theme.Current.TextColor

			// if the Layout is loaded...
			if loadState.loaded {
				// refer to the Layout's router
				router.Layout(gtx, th)
			} else {
				// otherwise, load the Layout into the loadstate
				loadState.Layout(gtx, th)
			}

			// and now draw the math to the screen
			e.Frame(gtx.Ops)
		}
	}
}

func main() {
	// start a go routine
	go func() {

		// run the app until lower levels err or exit
		err := runApp()

		if err != nil {
			log.Fatal(err)
		}
		// invite the os to exit the app
		os.Exit(0)
	}()

	//Main must be called last from the program main function.
	app.Main()

	/*
		On most platforms Main blocks forever,
		for Android and iOS it returns immediately
		to give control of the main thread back to the system.

		Calling Main is necessary because some operating systems
		require control of the main thread of the program for running windows.
	*/
}
