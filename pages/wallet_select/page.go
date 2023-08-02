package page_wallet_select

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/containers/bottom_bar"
	"github.com/g45t345rt/g45w/pages"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type Page struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation
	header         *prefabs.Header

	pageSelectWallet     *PageSelectWallet
	pageCreateWalletForm *PageCreateWalletForm

	pageRouter *router.Router
}

var _ router.Page = &Page{}

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

	pageRouter := router.NewRouter()

	pageSelectWallet := NewPageSelectWallet()
	pageRouter.Add(PAGE_SELECT_WALLET, pageSelectWallet)

	pageCreateWalletForm := NewPageCreateWalletForm()
	pageRouter.Add(PAGE_CREATE_WALLET_FORM, pageCreateWalletForm)

	pageCreateWalletSeedForm := NewPageCreateWalletSeedForm()
	pageRouter.Add(PAGE_CREATE_WALLET_SEED_FORM, pageCreateWalletSeedForm)

	pageCreateWalletHexSeedForm := NewPageCreateWalletHexSeedForm()
	pageRouter.Add(PAGE_CREATE_WALLET_HEXSEED_FORM, pageCreateWalletHexSeedForm)

	pageCreateWalletFastRegForm := NewPageCreateWalletFastRegForm()
	pageRouter.Add(PAGE_CREATE_WALLET_FASTREG_FORM, pageCreateWalletFastRegForm)

	pageCreateWalletDiskForm := NewPageCreateWalletDiskForm()
	pageRouter.Add(PAGE_CREATE_WALLET_DISK_FORM, pageCreateWalletDiskForm)

	header := prefabs.NewHeader(pageRouter)

	page := &Page{
		animationEnter:       animationEnter,
		animationLeave:       animationLeave,
		pageRouter:           pageRouter,
		header:               header,
		pageSelectWallet:     pageSelectWallet,
		pageCreateWalletForm: pageCreateWalletForm,
	}

	page_instance = page
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

	lastHistory := p.header.GetLastHistory()
	if lastHistory != nil {
		p.pageRouter.SetCurrent(lastHistory)
	} else {
		p.header.AddHistory(PAGE_SELECT_WALLET)
		p.pageRouter.SetCurrent(PAGE_SELECT_WALLET)
	}
}

func (p *Page) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if bottom_bar.Instance.ButtonWallet.Button.Clicked() {
		app_instance.Router.SetCurrent(pages.PAGE_WALLET)
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
				if state.Finished {
					p.isActive = false
					op.InvalidateOp{}.Add(gtx.Ops)
				}

				if state.Active {
					defer animation.TransformY(gtx, state.Value).Push(gtx.Ops).Pop()
				}
			}

			defer prefabs.PaintGrayLinearGradient(gtx).Pop()

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top: unit.Dp(30), Bottom: unit.Dp(30),
						Left: unit.Dp(30), Right: unit.Dp(30),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return p.header.Layout(gtx, th, func(gtx layout.Context, th *material.Theme, title string) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(22), title)
							lbl.Font.Weight = font.Bold

							return lbl.Layout(gtx)
						})
					})
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return p.pageRouter.Layout(gtx, th)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return bottom_bar.Instance.Layout(gtx, th)
		}),
	)
}
