package page_wallet_select

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/pages"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type Page struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation
	header         *pages.Header

	childRouter *router.Router
}

var _ router.Container = &Page{}

type PageInstance struct {
	router *router.Router
	header *pages.Header
}

var page_instance *PageInstance

func NewPage() *Page {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .5, ease.OutCubic),
	))

	childRouter := router.NewRouter()

	pageSelectWallet := NewPageSelectWallet()
	childRouter.Add("select_wallet", pageSelectWallet)

	pageCreateWalletForm := NewPageCreateWalletForm()
	childRouter.Add("create_wallet_form", pageCreateWalletForm)

	pageCreateWalletSeedForm := NewPageCreateWalletSeedForm()
	childRouter.Add("create_wallet_seed_form", pageCreateWalletSeedForm)

	pageCreateWalletHexSeedForm := NewPageCreateWalletHexSeedForm()
	childRouter.Add("create_wallet_hexseed_form", pageCreateWalletHexSeedForm)

	pageCreateWalletFastRegForm := NewPageCreateWalletFastRegForm()
	childRouter.Add("create_wallet_fastreg_form", pageCreateWalletFastRegForm)

	pageCreateWalletDiskForm := NewPageCreateWalletDiskForm()
	childRouter.Add("create_wallet_disk_form", pageCreateWalletDiskForm)

	th := app_instance.Current.Theme
	labelHeaderStyle := material.Label(th, unit.Sp(22), "")
	labelHeaderStyle.Font.Weight = font.Bold
	header := pages.NewHeader(labelHeaderStyle, childRouter)

	page_instance = &PageInstance{
		router: childRouter,
		header: header,
	}

	childRouter.SetPrimary("select_wallet")

	return &Page{
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		childRouter:    childRouter,
		header:         header,
	}
}

func (p *Page) IsActive() bool {
	return p.isActive
}

func (p *Page) Enter() {
	p.isActive = true
	p.animationLeave.Reset()
	p.animationEnter.Start()
}

func (p *Page) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	dr := image.Rectangle{Max: gtx.Constraints.Min}
	paint.LinearGradientOp{
		Stop1:  f32.Pt(0, float32(dr.Min.Y)),
		Stop2:  f32.Pt(0, float32(dr.Max.Y)),
		Color1: color.NRGBA{R: 0, G: 0, B: 0, A: 5},
		Color2: color.NRGBA{R: 0, G: 0, B: 0, A: 50},
	}.Add(gtx.Ops)
	defer clip.Rect(dr).Push(gtx.Ops).Pop()
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			{
				state := p.animationEnter.Update(gtx)
				if state.Active {
					animation.TransformY(gtx, state.Value).Add(gtx.Ops)
				}
			}

			{
				state := p.animationLeave.Update(gtx)
				if state.Finished {
					p.isActive = false
					op.InvalidateOp{}.Add(gtx.Ops)
				}

				if state.Active {
					animation.TransformY(gtx, state.Value).Add(gtx.Ops)
				}
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top: unit.Dp(30), Bottom: unit.Dp(30),
						Left: unit.Dp(30), Right: unit.Dp(30),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return p.header.Layout(gtx, th, nil)
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
