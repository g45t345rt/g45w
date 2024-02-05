package page_settings

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/containers/bottom_bar"
	"github.com/g45t345rt/g45w/pages"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
)

type Page struct {
	isActive bool

	pageSectionAnimation *pages.PageSectionAnimation
	header               *prefabs.Header
	pageRouter           *router.Router

	pageEditIPFSGateway *PageEditIPFSGateway

	pageMain    *PageMain
	pageAppInfo *PageAppInfo
}

var (
	PAGE_MAIN              = "page_main"
	PAGE_APP_INFO          = "page_app_info"
	PAGE_IPFS_GATEWAYS     = "page_ipfs_gateways"
	PAGE_ADD_IPFS_GATEWAY  = "page_add_ipfs_gateway"
	PAGE_EDIT_IPFS_GATEWAY = "page_edit_ipfs_gateway"
	PAGE_DONATION          = "page_donation"
)

var page_instance *Page

var _ router.Page = &Page{}

func New() *Page {
	pageRouter := router.NewRouter()

	pageMain := NewPageFront()
	pageRouter.Add(PAGE_MAIN, pageMain)

	pageAppInfo := NewPageAppInfo()
	pageRouter.Add(PAGE_APP_INFO, pageAppInfo)

	pageIPFSGateways := NewPageIPFSGateways()
	pageRouter.Add(PAGE_IPFS_GATEWAYS, pageIPFSGateways)

	pageAddIPFSGateway := NewPageAddIPFSGateway()
	pageRouter.Add(PAGE_ADD_IPFS_GATEWAY, pageAddIPFSGateway)

	pageEditIPFSGateway := NewPageEditIPFSGateway()
	pageRouter.Add(PAGE_EDIT_IPFS_GATEWAY, pageEditIPFSGateway)

	pageDonation := NewPageDonation()
	pageRouter.Add(PAGE_DONATION, pageDonation)

	header := prefabs.NewHeader(pageRouter)

	page := &Page{
		header:               header,
		pageRouter:           pageRouter,
		pageMain:             pageMain,
		pageAppInfo:          pageAppInfo,
		pageEditIPFSGateway:  pageEditIPFSGateway,
		pageSectionAnimation: pages.NewPageSectionAnimation(),
	}

	page_instance = page
	return page
}

func (p *Page) IsActive() bool {
	return p.isActive
}

func (p *Page) Enter() {
	bottom_bar.Instance.SetButtonActive(bottom_bar.BUTTON_SETTINGS)
	p.isActive = p.pageSectionAnimation.Enter()

	lastHistory := p.header.GetLastHistory()
	if lastHistory != nil {
		p.pageRouter.SetCurrent(lastHistory)
	} else {
		p.header.AddHistory(PAGE_MAIN)
		p.pageRouter.SetCurrent(PAGE_MAIN)
	}
}

func (p *Page) Leave() {
	p.isActive = p.pageSectionAnimation.Leave()
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			defer p.pageSectionAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

			startColor := theme.Current.BgGradientStartColor
			endColor := theme.Current.BgGradientEndColor
			defer utils.PaintLinearGradient(gtx, startColor, endColor).Pop()

			p.header.HandleKeyGoBack(gtx)
			p.header.HandleSwipeRightGoBack(gtx)

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top: theme.PagePadding, Bottom: theme.PagePadding,
						Left: theme.PagePadding, Right: theme.PagePadding,
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
