package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/camera"
	"github.com/deroproject/derohe/globals"
	"github.com/g45t345rt/g45w/android_background_service"
	"github.com/g45t345rt/g45w/app_db"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/bridge_metamask"
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

	// add android permissions
	_ "gioui.org/app/permission/camera"
	_ "gioui.org/app/permission/foreground_service"
	_ "gioui.org/app/permission/networkstate"
	_ "gioui.org/app/permission/post_notifications"
	_ "gioui.org/app/permission/storage"

	// support webp images
	_ "golang.org/x/image/webp"
)

// ANDROID_MANIFEST_SERVICE:org.gioui.x.worker_android$WorkerService

func loadFontCollection() ([]font.FontFace, error) {
	// universal fonts from https://github.com/satbyy/go-noto-universal

	goNotoKurrentRegularTTF, err := assets.GetFont("GoNotoKurrent-Regular.ttf")
	if err != nil {
		return nil, err
	}

	goNotoKurrentRegular, err := opentype.Parse(goNotoKurrentRegularTTF)
	if err != nil {
		return nil, err
	}

	goNotoKurrentBoldTTF, err := assets.GetFont("GoNotoKurrent-Bold.ttf")
	if err != nil {
		return nil, err
	}

	goNotoKurrentBold, err := opentype.Parse(goNotoKurrentBoldTTF)
	if err != nil {
		return nil, err
	}

	fontCollection := []font.FontFace{}
	fontCollection = append(fontCollection, font.FontFace{Font: font.Font{}, Face: goNotoKurrentRegular})
	fontCollection = append(fontCollection, font.FontFace{Font: font.Font{Weight: font.Bold}, Face: goNotoKurrentBold})
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

type TestnetTopBar struct{}

func (t TestnetTopBar) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if !settings.App.Testnet {
		return layout.Dimensions{}
	}

	return layout.N.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			r := op.Record(gtx.Ops)
			dims := layout.Inset{
				Left: unit.Dp(8), Right: unit.Dp(8),
				Top: unit.Dp(2), Bottom: unit.Dp(3),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(th, unit.Sp(16), lang.Translate("Testnet"))
				lbl.Font.Weight = font.Bold
				return lbl.Layout(gtx)
			})
			c := r.Stop()

			paint.FillShape(gtx.Ops, color.NRGBA{R: 255, A: 255}, clip.RRect{
				Rect: image.Rectangle{Max: dims.Size},
				SE:   gtx.Dp(5), SW: gtx.Dp(5),
			}.Op(gtx.Ops))

			c.Add(gtx.Ops)
			return dims

		})
	})
}

func loadDeroGlobals(useTestnet bool) {
	globals.Arguments["--testnet"] = useTestnet
	globals.Arguments["--debug"] = false
	globals.Arguments["--flog-level"] = nil
	globals.Arguments["--log-dir"] = nil
	globals.Arguments["--help"] = false
	globals.Arguments["--version"] = false
	globals.InitNetwork() // this func assign mainnet/testnet config depending on globals.Arguments["--testnet"] value
}

func runApp() error {
	var ops op.Ops
	app_instance.Load()
	window := app_instance.Window
	explorer := app_instance.Explorer
	router := app_instance.Router

	theme.LoadImages()
	loadState := NewLoadState(window)

	go bridge_metamask.StartServer()

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

		if android_background_service.IsAvailable() {
			err = android_background_service.Start()
			if err != nil {
				loadState.SetStatus("", err)
				return
			}

			if settings.App.MobileBackgroundService {
				android_background_service.StartForeground()
			}
		}

		loadDeroGlobals(settings.App.Testnet)

		loadState.SetStatus(lang.Translate("Loading lookup table"), nil)
		//walletapi.Initialize_LookupTable(1, 1<<21)
		err = lookup_table.Load()
		if err != nil {
			loadState.SetStatus("", err)
			return
		}

		loadState.SetStatus(lang.Translate("Loading app data"), nil)
		err = app_db.Load()
		if err != nil {
			loadState.SetStatus("", err)
			return
		}

		/*
			loadState.SetStatus(lang.Translate("Loading wallets"), nil)
			err = wallet_manager.Load()
			if err != nil {
				loadState.SetStatus("", err)
				return
			}*/

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

	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(fontCollection))
	th.FingerSize = 48

	for {
		e := window.NextEvent()
		explorer.ListenEvents(e)
		camera.ListenEvents(e)
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			paint.Fill(gtx.Ops, theme.Current.BgColor)
			th.Bg = theme.Current.BgColor
			th.Fg = theme.Current.TextColor

			if loadState.loaded {
				router.Layout(gtx, th)
			} else {
				loadState.Layout(gtx, th)
			}

			TestnetTopBar{}.Layout(gtx, th)
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
