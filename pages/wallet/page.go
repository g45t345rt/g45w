package page_wallet

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/pages"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
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

	header         *prefabs.Header
	buttonCopyAddr *components.Button

	pageBalanceTokens *PageBalanceTokens
	pageSendForm      *PageSendForm

	childRouter *router.Router
	infoModal   *components.NotificationModal
}

var _ router.Container = &Page{}

type PageInstance struct {
	router *router.Router
	header *prefabs.Header
}

var page_instance *PageInstance

func NewPage() *Page {
	th := app_instance.Current.Theme

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .5, ease.OutCubic),
	))

	childRouter := router.NewRouter()
	pageBalanceTokens := NewPageBalanceTokens()
	childRouter.Add("balanceTokens", pageBalanceTokens)
	pageSendForm := NewPageSendForm()
	childRouter.Add("sendForm", pageSendForm)
	pageReceiveForm := NewPageReceiveForm()
	childRouter.Add("receiveForm", pageReceiveForm)
	pageSettings := NewPageSettings()
	childRouter.Add("settings", pageSettings)
	pageAddSCForm := NewPageAddSCForm()
	childRouter.Add("addSCForm", pageAddSCForm)

	labelHeaderStyle := material.Label(th, unit.Sp(22), "")
	labelHeaderStyle.Font.Weight = font.Bold

	settingsIcon, _ := widget.NewIcon(icons.ActionSettings)
	buttonSettings := components.NewButton(components.ButtonStyle{
		Icon:      settingsIcon,
		TextColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
	})

	header := prefabs.NewHeader(labelHeaderStyle, childRouter, buttonSettings)

	page_instance = &PageInstance{
		router: childRouter,
		header: header,
	}

	childRouter.SetPrimary("balanceTokens")

	textColor := color.NRGBA{R: 0, G: 0, B: 0, A: 100}
	textHoverColor := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

	copyIcon, _ := widget.NewIcon(icons.ContentContentCopy)
	buttonCopyAddr := components.NewButton(components.ButtonStyle{
		Icon:           copyIcon,
		TextColor:      textColor,
		HoverTextColor: &textHoverColor,
	})

	w := app_instance.Current.Window
	infoModal := components.NewNotificationInfoModal(w)

	router := app_instance.Current.Router
	router.PushLayout(func(gtx layout.Context, th *material.Theme) {
		infoModal.Layout(gtx, th)
	})

	return &Page{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		header: header,

		buttonCopyAddr:    buttonCopyAddr,
		pageBalanceTokens: pageBalanceTokens,
		pageSendForm:      pageSendForm,
		childRouter:       childRouter,
		infoModal:         infoModal,
	}
}

func (p *Page) SetCurrent(tag string) {
	p.childRouter.SetCurrent(tag)
}

func (p *Page) IsActive() bool {
	return p.isActive
}

func (p *Page) Enter() {
	openedWallet := wallet_manager.Instance.OpenedWallet
	if openedWallet != nil {
		p.isActive = true
		pages.BottomBarInstance.SetButtonActive("wallet")
		p.header.SetTitle(fmt.Sprintf("Wallet [%s]", openedWallet.Info.Name))

		w := app_instance.Current.Window
		w.Option(app.StatusColor(color.NRGBA{A: 255}))

		p.animationLeave.Reset()
		p.animationEnter.Start()
	} else {
		app_instance.Current.Router.SetCurrent("page_wallet_select")
	}
}

func (p *Page) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	openedWallet := wallet_manager.Instance.OpenedWallet
	if openedWallet == nil {
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}

	walletAddr := utils.ReduceAddr(openedWallet.Info.Addr)

	if p.pageBalanceTokens.displayBalance.buttonSend.Clickable.Clicked() {
		p.childRouter.SetCurrent("sendForm")
	}

	if p.pageBalanceTokens.displayBalance.buttonReceive.Clickable.Clicked() {
		p.childRouter.SetCurrent("receiveForm")
	}

	if p.pageBalanceTokens.tokenBar.buttonAddToken.Clickable.Clicked() {
		p.childRouter.SetCurrent("addSCForm")
	}

	if p.header.ButtonRight.Clickable.Clicked() {
		p.childRouter.SetCurrent("settings")
	}

	if p.buttonCopyAddr.Clickable.Clicked() {
		clipboard.WriteOp{
			Text: walletAddr,
		}.Add(gtx.Ops)
		p.infoModal.SetText("Clipboard", "Addr copied to clipboard")
		p.infoModal.SetVisible(true)
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

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				//layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				//	return app_instance.Current.NodeStatusBar.Layout(gtx, th)
				//}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.RRect{
						Rect: image.Rectangle{Max: gtx.Constraints.Max},
						NW:   gtx.Dp(15),
						NE:   gtx.Dp(15),
						SE:   gtx.Dp(0),
						SW:   gtx.Dp(0),
					}.Op(gtx.Ops))

					dr := image.Rectangle{Max: gtx.Constraints.Max}
					paint.LinearGradientOp{
						Stop1:  f32.Pt(0, float32(dr.Min.Y)),
						Stop2:  f32.Pt(0, float32(dr.Max.Y)),
						Color1: color.NRGBA{R: 0, G: 0, B: 0, A: 5},
						Color2: color.NRGBA{R: 0, G: 0, B: 0, A: 50},
					}.Add(gtx.Ops)
					defer clip.Rect(dr).Push(gtx.Ops).Pop()
					paint.PaintOp{}.Add(gtx.Ops)

					return layout.Inset{
						Left: unit.Dp(30), Right: unit.Dp(30),
						Top: unit.Dp(30), Bottom: unit.Dp(20),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return p.header.Layout(gtx, th, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									label := material.Label(th, unit.Sp(16), walletAddr)
									label.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 200}
									return label.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Width: unit.Dp(5)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									gtx.Constraints.Max.X = gtx.Dp(18)
									gtx.Constraints.Max.Y = gtx.Dp(18)
									return p.buttonCopyAddr.Layout(gtx, th)
								}),
							)
						})
					})
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return p.childRouter.Layout(gtx, th)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return pages.BottomBarInstance.Layout(gtx, th)
		}),
	)
}
