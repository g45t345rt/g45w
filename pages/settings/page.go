package page_settings

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/containers/bottom_bar"
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
	pageRouter     *router.Router

	pageEditIPFSGateway *PageEditIPFSGateway

	pageMain *PageMain
	pageInfo *PageInfo
}

var (
	PAGE_MAIN              = "page_main"
	PAGE_INFO              = "page_info"
	PAGE_IPFS_GATEWAYS     = "page_ipfs_gateways"
	PAGE_ADD_IPFS_GATEWAY  = "page_add_ipfs_gateway"
	PAGE_EDIT_IPFS_GATEWAY = "page_edit_ipfs_gateway"
)

var page_instance *Page

var _ router.Page = &Page{}

func New() *Page {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .5, ease.OutCubic),
	))

	pageRouter := router.NewRouter()

	pageMain := NewPageFront()
	pageRouter.Add(PAGE_MAIN, pageMain)

	pageInfo := NewPageInfo()
	pageRouter.Add(PAGE_INFO, pageInfo)

	pageIPFSGateways := NewPageIPFSGateways()
	pageRouter.Add(PAGE_IPFS_GATEWAYS, pageIPFSGateways)

	pageAddIPFSGateway := NewPageAddIPFSGateway()
	pageRouter.Add(PAGE_ADD_IPFS_GATEWAY, pageAddIPFSGateway)

	pageEditIPFSGateway := NewPageEditIPFSGateway()
	pageRouter.Add(PAGE_EDIT_IPFS_GATEWAY, pageEditIPFSGateway)

	header := prefabs.NewHeader(pageRouter)

	page := &Page{
		animationEnter:      animationEnter,
		animationLeave:      animationLeave,
		header:              header,
		pageRouter:          pageRouter,
		pageMain:            pageMain,
		pageInfo:            pageInfo,
		pageEditIPFSGateway: pageEditIPFSGateway,
	}

	page_instance = page
	return page
}

func (p *Page) IsActive() bool {
	return p.isActive
}

func (p *Page) Enter() {
	bottom_bar.Instance.SetButtonActive(bottom_bar.BUTTON_SETTINGS)
	p.animationEnter.Start()
	p.animationLeave.Reset()

	lastHistory := p.header.GetLastHistory()
	if lastHistory != nil {
		p.pageRouter.SetCurrent(lastHistory)
	} else {
		p.header.AddHistory(PAGE_MAIN)
		p.pageRouter.SetCurrent(PAGE_MAIN)
	}

	p.isActive = true
}

func (p *Page) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

			defer prefabs.PaintLinearGradient(gtx).Pop()

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
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.pageRouter.Layout(gtx, th)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return bottom_bar.Instance.Layout(gtx, th)
		}),
	)
}
