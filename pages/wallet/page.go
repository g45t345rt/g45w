package page_wallet

import (
	"fmt"
	"image"
	"image/color"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/creachadair/jrpc2"
	"github.com/deroproject/derohe/walletapi/xswd"
	"github.com/g45t345rt/g45w/android_background_service"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/containers/bottom_bar"
	"github.com/g45t345rt/g45w/containers/confirm_modal"
	"github.com/g45t345rt/g45w/containers/node_status_bar"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/containers/xswd_perm_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/pages"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Page struct {
	isActive       bool
	animationEnter *animation.Animation
	animationLeave *animation.Animation

	header     *prefabs.Header
	xswdHeader *XSWDHeader

	pageBalanceTokens   *PageBalanceTokens
	pageSendForm        *PageSendForm
	pageSCToken         *PageSCToken
	pageContactForm     *PageContactForm
	pageSendOptionsForm *PageSendOptionsForm
	pageSCFolders       *PageSCFolders
	pageContacts        *PageContacts
	pageTransaction     *PageTransaction
	pageDexSwap         *PageDEXSwap
	pageDEXAddLiquidity *PageDEXAddLiquidity
	pageDEXRemLiquidity *PageDEXRemLiquidity
	pageDEXSCBridgeOut  *PageDEXSCBridgeOut
	pageDEXSCBridgeIn   *PageDEXSCBridgeIn
	pageXSWDManage      *PageXSWDManage
	pageXSWDApp         *PageXSWDApp
	pageSCExplorer      *PageSCExplorer

	pageRouter *router.Router
}

var _ router.Page = &Page{}

var page_instance *Page

var (
	PAGE_SETTINGS          = "page_settings"
	PAGE_SEND_FORM         = "page_send_form"
	PAGE_RECEIVE_FORM      = "page_receive_form"
	PAGE_BALANCE_TOKENS    = "page_balance_tokens"
	PAGE_ADD_SC_FORM       = "page_add_sc_form"
	PAGE_TXS               = "page_txs"
	PAGE_SC_TOKEN          = "page_sc_token"
	PAGE_REGISTER_WALLET   = "page_register_wallet"
	PAGE_CONTACTS          = "page_contacts"
	PAGE_CONTACT_FORM      = "page_contact_form"
	PAGE_SEND_OPTIONS_FORM = "page_send_options_form"
	PAGE_SC_FOLDERS        = "page_sc_folders"
	PAGE_WALLET_INFO       = "page_wallet_info"
	PAGE_TRANSACTION       = "page_transaction"
	PAGE_SCAN_COLLECTION   = "page_scan_collection"
	PAGE_SERVICE_NAMES     = "page_service_names"
	PAGE_DEX_PAIRS         = "page_dex_pairs"
	PAGE_DEX_SWAP          = "page_dex_swap"
	PAGE_DEX_ADD_LIQUIDITY = "page_dex_add_liquidity"
	PAGE_DEX_REM_LIQUIDITY = "page_dex_rem_liquidity"
	PAGE_DEX_SC_BRIDGE_OUT = "page_dex_sc_bridge_out"
	PAGE_DEX_SC_BRIDGE_IN  = "page_dex_sc_bridge_in"
	PAGE_XSWD_MANAGE       = "page_xswd_manage"
	PAGE_XSWD_APP          = "page_xswd_app"
	PAGE_SC_EXPLORER       = "page_sc_explorer"
)

func New() *Page {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .5, ease.OutCubic),
	))

	pageRouter := router.NewRouter()
	pageBalanceTokens := NewPageBalanceTokens()
	pageRouter.Add(PAGE_BALANCE_TOKENS, pageBalanceTokens)

	pageSendForm := NewPageSendForm()
	pageRouter.Add(PAGE_SEND_FORM, pageSendForm)

	pageReceiveForm := NewPageReceiveForm()
	pageRouter.Add(PAGE_RECEIVE_FORM, pageReceiveForm)

	pageSettings := NewPageSettings()
	pageRouter.Add(PAGE_SETTINGS, pageSettings)

	pageAddSCForm := NewPageAddSCForm()
	pageRouter.Add(PAGE_ADD_SC_FORM, pageAddSCForm)

	pageSCToken := NewPageSCToken()
	pageRouter.Add(PAGE_SC_TOKEN, pageSCToken)

	pageRegisterWallet := NewPageRegisterWallet()
	pageRouter.Add(PAGE_REGISTER_WALLET, pageRegisterWallet)

	pageContacts := NewPageContacts()
	pageRouter.Add(PAGE_CONTACTS, pageContacts)

	pageContactForm := NewPageContactForm()
	pageRouter.Add(PAGE_CONTACT_FORM, pageContactForm)

	pageSendOptionsForm := NewPageSendOptionsForm()
	pageRouter.Add(PAGE_SEND_OPTIONS_FORM, pageSendOptionsForm)

	pageSCFolders := NewPageSCFolders()
	pageRouter.Add(PAGE_SC_FOLDERS, pageSCFolders)

	pageWalletInfo := NewPageWalletInfo()
	pageRouter.Add(PAGE_WALLET_INFO, pageWalletInfo)

	pageTransaction := NewPageTransaction()
	pageRouter.Add(PAGE_TRANSACTION, pageTransaction)

	pageScanCollection := NewPageScanCollection()
	pageRouter.Add(PAGE_SCAN_COLLECTION, pageScanCollection)

	pageServiceNames := NewPageServiceNames()
	pageRouter.Add(PAGE_SERVICE_NAMES, pageServiceNames)

	pageDEXPairs := NewPageDEXPairs()
	pageRouter.Add(PAGE_DEX_PAIRS, pageDEXPairs)

	pageDEXSwap := NewPageDEXSwap()
	pageRouter.Add(PAGE_DEX_SWAP, pageDEXSwap)

	pageDEXAddLiquidity := NewPageDEXAddLiquidity()
	pageRouter.Add(PAGE_DEX_ADD_LIQUIDITY, pageDEXAddLiquidity)

	pageDEXRemLiquidity := NewPageDEXRemLiquidity()
	pageRouter.Add(PAGE_DEX_REM_LIQUIDITY, pageDEXRemLiquidity)

	pageDEXSCBridgeOut := NewPageDEXSCBridgeOut()
	pageRouter.Add(PAGE_DEX_SC_BRIDGE_OUT, pageDEXSCBridgeOut)

	pageDEXSCBridgeIn := NewPageDEXSCBridgeIn()
	pageRouter.Add(PAGE_DEX_SC_BRIDGE_IN, pageDEXSCBridgeIn)

	pageXSWDManage := NewPageXSWDManage()
	pageRouter.Add(PAGE_XSWD_MANAGE, pageXSWDManage)

	pageXSWDApp := NewPageXSWDApp()
	pageRouter.Add(PAGE_XSWD_APP, pageXSWDApp)

	pageSCExplorer := NewPageSCExplorer()
	pageRouter.Add(PAGE_SC_EXPLORER, pageSCExplorer)

	header := prefabs.NewHeader(pageRouter)

	page := &Page{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		header:     header,
		xswdHeader: NewXSWDHeader(),

		pageBalanceTokens:   pageBalanceTokens,
		pageSendForm:        pageSendForm,
		pageSCToken:         pageSCToken,
		pageContactForm:     pageContactForm,
		pageSendOptionsForm: pageSendOptionsForm,
		pageSCFolders:       pageSCFolders,
		pageContacts:        pageContacts,
		pageTransaction:     pageTransaction,
		pageDexSwap:         pageDEXSwap,
		pageDEXAddLiquidity: pageDEXAddLiquidity,
		pageDEXRemLiquidity: pageDEXRemLiquidity,
		pageDEXSCBridgeOut:  pageDEXSCBridgeOut,
		pageDEXSCBridgeIn:   pageDEXSCBridgeIn,
		pageXSWDManage:      pageXSWDManage,
		pageXSWDApp:         pageXSWDApp,
		pageSCExplorer:      pageSCExplorer,

		pageRouter: pageRouter,
	}
	page_instance = page
	return page
}

func (p *Page) IsActive() bool {
	return p.isActive
}

func (p *Page) Enter() {
	bottom_bar.Instance.SetButtonActive(bottom_bar.BUTTON_WALLET)
	openedWallet := wallet_manager.OpenedWallet
	if openedWallet != nil {
		p.isActive = true
		w := app_instance.Window
		w.Option(app.StatusColor(color.NRGBA{A: 255}))

		p.animationLeave.Reset()
		p.animationEnter.Start()

		lastHistory := p.header.GetLastHistory()
		if lastHistory != nil {
			p.pageRouter.SetCurrent(lastHistory)
		} else {
			p.header.AddHistory(PAGE_BALANCE_TOKENS)
			p.pageRouter.SetCurrent(PAGE_BALANCE_TOKENS)
		}

		p.askToCreateFolderTokens()
	} else {
		app_instance.Router.SetCurrent(pages.PAGE_WALLET_SELECT)
	}
}

func (p *Page) CloseXSWD() error {
	wallet := wallet_manager.OpenedWallet
	w := app_instance.Window

	if android_background_service.IsAvailable() {
		if !settings.App.MobileBackgroundService {
			running, err := android_background_service.IsRunning()
			if err != nil {
				return err
			}

			if running {
				err := android_background_service.Stop()
				if err != nil {
					return err
				}
			}
		}
	}

	wallet.CloseXSWD()
	w.Invalidate()
	return nil
}

func (p *Page) OpenXSWD() error {
	wallet := wallet_manager.OpenedWallet
	w := app_instance.Window

	if android_background_service.IsAvailable() {
		running, err := android_background_service.IsRunning()
		if err != nil {
			return err
		}

		if !running {
			if wallet.Settings.NotifyXSWDMobileBackgroundService {
				yes := <-confirm_modal.Instance.Open(confirm_modal.ConfirmText{
					Title:  lang.Translate("Background Service"),
					Prompt: lang.Translate("Using XSWD on mobile requires running the application in the background."),
					Yes:    lang.Translate("OK"),
					No:     lang.Translate("CANCEL"),
				})

				if !yes {
					return nil
				}

				wallet.Settings.NotifyXSWDMobileBackgroundService = false
				err := wallet.SaveSettings()
				if err != nil {
					return err
				}
			}

			err = android_background_service.Start()
			if err != nil {
				return err
			}
		}
	}

	// TODO: IOS background app with GPS hack

	appHandler := func(appData *xswd.ApplicationData) bool {
		prompt := lang.Translate("The app [{}] is trying to connect. Do you want to give permission?")
		prompt = strings.Replace(prompt, "{}", appData.Name, -1)
		yesChan := confirm_modal.Instance.Open(confirm_modal.ConfirmText{
			Title:  lang.Translate("XSWD Auth"),
			Prompt: prompt,
		})

		w.Invalidate()
		yes := <-yesChan
		go func() {
			time.Sleep(100 * time.Millisecond)
			page_instance.pageXSWDManage.Load()
			w.Invalidate()
		}()

		return yes
	}

	reqHandler := func(appData *xswd.ApplicationData, req *jrpc2.Request) xswd.Permission {
		fmt.Println(req.Method(), req.ParamString())
		permChan := xswd_perm_modal.Instance.Open(appData, req.Method())
		w.Invalidate()
		return <-permChan
	}

	wallet.OpenXSWD(appHandler, reqHandler)
	w.Invalidate()
	return nil
}

func (p *Page) askToCreateFolderTokens() {
	wallet := wallet_manager.OpenedWallet

	// make you don't ask if testnet mode or was already asked before
	if !wallet.Settings.AskToStoreDEXTokens ||
		settings.App.Testnet {
		return
	}

	go func() {
		yes := <-confirm_modal.Instance.Open(confirm_modal.ConfirmText{
			Title:  lang.Translate("DEX Tokens"),
			Prompt: lang.Translate("Would you like to create a folder that includes all DEX tokens? Popular tokens will be automatically added to your favorites."),
		})

		if yes {
			err := wallet.InsertDexTokensFolder()
			if err != nil {
				notification_modal.Open(notification_modal.Params{
					Type:  notification_modal.ERROR,
					Title: lang.Translate("Error"),
					Text:  err.Error(),
				})
			} else {
				notification_modal.Open(notification_modal.Params{
					Type:  notification_modal.SUCCESS,
					Title: lang.Translate("Success"),
					Text:  lang.Translate("The folder was created."),
				})
			}
		}

		wallet.Settings.AskToStoreDEXTokens = false
		err := wallet.SaveSettings()
		if err != nil {
			notification_modal.Open(notification_modal.Params{
				Type:  notification_modal.ERROR,
				Title: lang.Translate("Error"),
				Text:  err.Error(),
			})
		}

		page_instance.pageBalanceTokens.LoadTokens()
	}()
}

func (p *Page) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	openedWallet := wallet_manager.OpenedWallet
	if openedWallet == nil {
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			{
				state := p.animationEnter.Update(gtx)
				if state.Active {
					defer animation.TransformY(gtx, state.Value).Push(gtx.Ops).Pop()
				}
			}

			{
				state := p.animationLeave.Update(gtx)

				if state.Active {
					defer animation.TransformY(gtx, state.Value).Push(gtx.Ops).Pop()
				}

				if state.Finished {
					p.isActive = false
					op.InvalidateOp{}.Add(gtx.Ops)
				}
			}

			p.header.HandleKeyGoBack(gtx)
			p.header.HandleSwipeRightGoBack(gtx)

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return node_status_bar.Instance.Layout(gtx, th)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					r := op.Record(gtx.Ops)
					dims := layout.Inset{
						Left: theme.PagePadding, Right: theme.PagePadding,
						Top: theme.PagePadding, Bottom: theme.PagePadding,
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						dims := p.header.Layout(gtx, th, func(gtx layout.Context, th *material.Theme, title string) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(22), title)
							lbl.Font.Weight = font.Bold
							return lbl.Layout(gtx)
						})

						if page_instance.pageRouter.Current == PAGE_BALANCE_TOKENS {
							p.xswdHeader.Layout(gtx, th)
						}

						return dims
					})
					c := r.Stop()

					paint.FillShape(gtx.Ops, theme.Current.HeaderTopBgColor, clip.RRect{
						Rect: image.Rectangle{Max: dims.Size},
						NW:   gtx.Dp(15),
						NE:   gtx.Dp(15),
						SE:   gtx.Dp(0),
						SW:   gtx.Dp(0),
					}.Op(gtx.Ops))

					c.Add(gtx.Ops)
					return dims
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					startColor := theme.Current.BgGradientStartColor
					endColor := theme.Current.BgGradientEndColor
					defer utils.PaintLinearGradient(gtx, startColor, endColor).Pop()

					return p.pageRouter.Layout(gtx, th)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return bottom_bar.Instance.Layout(gtx, th)
		}),
	)
}

type XSWDHeader struct {
	toggleClickable *widget.Clickable
	manageClickable *widget.Clickable
	manageIcon      *widget.Icon
}

func NewXSWDHeader() *XSWDHeader {
	manageIcon, _ := widget.NewIcon(icons.ActionSettingsApplications)

	return &XSWDHeader{
		toggleClickable: new(widget.Clickable),
		manageClickable: new(widget.Clickable),
		manageIcon:      manageIcon,
	}
}

func (x *XSWDHeader) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	wallet := wallet_manager.OpenedWallet
	xswd := wallet.ServerXSWD

	r := op.Record(gtx.Ops)
	width := gtx.Dp(20) // 20 because of left + right border and then we calculate each layout.Rigid
	dims := layout.Inset{
		Top: unit.Dp(5), Bottom: unit.Dp(5),
		Left: unit.Dp(10), Right: unit.Dp(10),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if x.toggleClickable.Clicked(gtx) {
					go func() {
						if xswd != nil && xswd.IsRunning() {
							page_instance.CloseXSWD()
						} else {
							err := page_instance.OpenXSWD()
							if err != nil {
								notification_modal.Open(notification_modal.Params{
									Type:  notification_modal.ERROR,
									Title: lang.Translate("Error"),
									Text:  err.Error(),
								})
							}
						}
					}()
				}

				dims := x.toggleClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					if x.toggleClickable.Hovered() {
						pointer.CursorPointer.Add(gtx.Ops)
					}

					txt := ""
					if xswd != nil && xswd.IsRunning() {
						txt = lang.Translate("XSWD ON")
					} else {
						txt = lang.Translate("XSWD OFF")
					}

					lbl := material.Label(th, unit.Sp(14), txt)
					lbl.Color = theme.Current.XSWDBgTextColor
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				})

				width += dims.Size.X
				return dims
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				dims := layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
				width += dims.Size.X
				return dims
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if x.manageClickable.Clicked(gtx) {
					page_instance.pageRouter.SetCurrent(PAGE_XSWD_MANAGE)
					page_instance.header.AddHistory(PAGE_XSWD_MANAGE)
				}

				dims := x.manageClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					if x.manageClickable.Hovered() {
						pointer.CursorPointer.Add(gtx.Ops)
					}

					gtx.Constraints.Min.X = gtx.Dp(20)
					gtx.Constraints.Min.Y = gtx.Dp(20)
					return x.manageIcon.Layout(gtx, theme.Current.XSWDBgTextColor)
				})
				width += dims.Size.X
				return dims
			}),
		)
	})
	c := r.Stop()

	// important -> dims.Size.X != width because we use layout.Flex{} and it keeps the full width

	offset := f32.Affine2D{}.Offset(f32.Pt(float32((dims.Size.X/2)-(width/2)), -float32(gtx.Dp(theme.PagePadding+10))))
	defer op.Affine(offset).Push(gtx.Ops).Pop()

	paint.FillShape(gtx.Ops, theme.Current.XSWDBgColor, clip.RRect{
		Rect: image.Rectangle{Max: image.Pt(int(width), dims.Size.Y)},
		NE:   gtx.Dp(10), NW: gtx.Dp(10),
		SE: gtx.Dp(10), SW: gtx.Dp(10),
	}.Op(gtx.Ops))
	c.Add(gtx.Ops)
	return dims
}
