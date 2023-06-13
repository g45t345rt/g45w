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
	"github.com/g45t345rt/g45w/containers/bottom_bar"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type Page struct {
	isActive bool

	animationEnter   *animation.Animation
	animationLeave   *animation.Animation
	header           *prefabs.Header
	pageSelectWallet *PageSelectWallet

	childRouter *router.Router
}

var _ router.Container = &Page{}

var page_instance *Page

const (
	PAGE_CREATE_WALLET_SEED_FORM    = "page_create_wallet_seed_form"
	PAGE_CREATE_WALLET_HEXSEED_FORM = "page_create_wallet_hexseed_form"
	PAGE_CREATE_WALLET_FORM         = "page_create_wallet_form"
	PAGE_CREATE_WALLET_FASTREG_FORM = "page_create_Wallet_fastreg_form"
	PAGE_CREATE_WALLET_DISK_FORM    = "page_create_wallet_disk_form"
	PAGE_SELECT_WALLET              = "page_select_wallet"
)

func New() *Page {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .5, ease.OutCubic),
	))

	childRouter := router.NewRouter()

	pageSelectWallet := NewPageSelectWallet()
	childRouter.Add(PAGE_SELECT_WALLET, pageSelectWallet)

	pageCreateWalletForm := NewPageCreateWalletForm()
	childRouter.Add(PAGE_CREATE_WALLET_FORM, pageCreateWalletForm)

	pageCreateWalletSeedForm := NewPageCreateWalletSeedForm()
	childRouter.Add(PAGE_CREATE_WALLET_SEED_FORM, pageCreateWalletSeedForm)

	pageCreateWalletHexSeedForm := NewPageCreateWalletHexSeedForm()
	childRouter.Add(PAGE_CREATE_WALLET_HEXSEED_FORM, pageCreateWalletHexSeedForm)

	pageCreateWalletFastRegForm := NewPageCreateWalletFastRegForm()
	childRouter.Add(PAGE_CREATE_WALLET_FASTREG_FORM, pageCreateWalletFastRegForm)

	pageCreateWalletDiskForm := NewPageCreateWalletDiskForm()
	childRouter.Add(PAGE_CREATE_WALLET_DISK_FORM, pageCreateWalletDiskForm)

	th := app_instance.Theme
	labelHeaderStyle := material.Label(th, unit.Sp(22), "")
	labelHeaderStyle.Font.Weight = font.Bold
	header := prefabs.NewHeader(labelHeaderStyle, childRouter, nil)

	page := &Page{
		animationEnter:   animationEnter,
		animationLeave:   animationLeave,
		childRouter:      childRouter,
		header:           header,
		pageSelectWallet: pageSelectWallet,
	}

	page_instance = page
	childRouter.SetPrimary(PAGE_SELECT_WALLET)

	return page
}

func (p *Page) IsActive() bool {
	return p.isActive
}

func (p *Page) Enter() {
	bottom_bar.Instance.SetButtonActive(bottom_bar.BUTTON_WALLET)
	p.isActive = true
	p.animationLeave.Reset()
	p.animationEnter.Start()

	p.pageSelectWallet.Load()
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
			return bottom_bar.Instance.Layout(gtx, th)
		}),
	)
}
