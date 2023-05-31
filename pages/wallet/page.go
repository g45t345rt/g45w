package page_wallet

import (
	"fmt"
	"image"
	"image/color"

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
	"github.com/deroproject/derohe/p2p"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/node"
	"github.com/g45t345rt/g45w/pages"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Page struct {
	isActive       bool
	animationEnter *animation.Animation
	animationLeave *animation.Animation

	nodeStatusBar  *NodeStatusBar
	header         *pages.Header
	buttonCopyAddr *components.Button

	pageBalanceTokens *PageBalanceTokens
	pageSendForm      *PageSendForm

	childRouter *router.Router
}

var _ router.Container = &Page{}

type PageInstance struct {
	router *router.Router
	header *pages.Header
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

	labelHeaderStyle := material.Label(th, unit.Sp(22), "")
	labelHeaderStyle.Font.Weight = font.Bold
	header := pages.NewHeader(labelHeaderStyle, childRouter)

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

	return &Page{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		nodeStatusBar: NewNodeStatusBar(),
		header:        header,

		buttonCopyAddr:    buttonCopyAddr,
		pageBalanceTokens: pageBalanceTokens,
		pageSendForm:      pageSendForm,
		childRouter:       childRouter,
	}
}

func (p *Page) SetCurrent(tag string) {
	p.childRouter.SetCurrent(tag)
}

func (p *Page) IsActive() bool {
	return p.isActive
}

func (p *Page) Enter() {
	p.isActive = true
	app_instance.Current.BottomBar.SetActive("wallet")
	p.header.LabelTitle.Text = "Wallet 0"
	//p.header.WalletAddr = "derog54...g435450"

	p.animationLeave.Reset()
	p.animationEnter.Start()
}

func (p *Page) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if p.pageBalanceTokens.displayBalance.buttonSend.Clickable.Clicked() {
		p.childRouter.SetCurrent("sendForm")
	}

	if p.pageBalanceTokens.displayBalance.buttonReceive.Clickable.Clicked() {
		p.childRouter.SetCurrent("receiveForm")
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
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.nodeStatusBar.Layout(gtx, th)
				}),
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
									label := material.Label(th, unit.Sp(16), "derog54...g435450")
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
			return app_instance.Current.BottomBar.Layout(gtx, th)
		}),
	)
}

type NodeStatusBar struct {
	clickable *widget.Clickable
}

func NewNodeStatusBar() *NodeStatusBar {
	return &NodeStatusBar{
		clickable: new(widget.Clickable),
	}
}

func (n *NodeStatusBar) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	paint.FillShape(gtx.Ops, color.NRGBA{A: 255}, clip.Rect{
		Max: gtx.Constraints.Max,
	}.Op())

	//paint.ColorOp{Color: color.NRGBA{A: 255}}.Add(gtx.Ops)
	//paint.PaintOp{}.Add(gtx.Ops)

	chain := node.Instance.Chain
	our_height := chain.Get_Height()
	best_height, _ := p2p.Best_Peer_Height()

	if n.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
	}

	if n.clickable.Clicked() {
		app_instance.Current.Router.SetCurrent("page_node")
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	return n.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(15), Bottom: unit.Dp(15),
			Left: unit.Dp(10), Right: unit.Dp(10),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = gtx.Dp(12)
					gtx.Constraints.Max.Y = gtx.Dp(12)
					paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 0, B: 0, A: 255},
						clip.Ellipse{
							Max: gtx.Constraints.Max,
						}.Op(gtx.Ops))

					return layout.Dimensions{Size: gtx.Constraints.Max}
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					status := fmt.Sprintf("%d / %d", our_height, best_height)
					label := material.Label(th, unit.Sp(16), status)
					label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
					return label.Layout(gtx)
				}),
			)
		})
	})
}
