package page_node

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/containers/bottom_bar"
	"github.com/g45t345rt/g45w/node_manager"
	"github.com/g45t345rt/g45w/pages"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
)

type Page struct {
	isActive bool

	pageRouter *router.Router

	pageSelectNode     *PageSelectNode
	pageAddNodeForm    *PageAddNodeForm
	pageEditNodeForm   *PageEditNodeForm
	pageRemoteNode     *PageRemoteNode
	pageIntegratedNode *PageIntegratedNode
	header             *prefabs.Header

	pageSectionAnimation *pages.PageSectionAnimation
}

var _ router.Page = &Page{}

var page_instance *Page

const (
	PAGE_SELECT_NODE     = "page_select_node"
	PAGE_ADD_NODE_FORM   = "page_add_node_form"
	PAGE_EDIT_NODE_FORM  = "page_edit_node_form"
	PAGE_INTEGRATED_NODE = "page_integrated_node"
	PAGE_REMOTE_NODE     = "page_remote_node"
)

func New() *Page {
	pageRouter := router.NewRouter()
	pageSelectNode := NewPageSelectNode()
	pageRouter.Add(PAGE_SELECT_NODE, pageSelectNode)

	pageAddNodeForm := NewPageAddNodeForm()
	pageRouter.Add(PAGE_ADD_NODE_FORM, pageAddNodeForm)

	pageEditNodeForm := NewPageEditNodeForm()
	pageRouter.Add(PAGE_EDIT_NODE_FORM, pageEditNodeForm)

	pageIntegratedNode := NewPageIntegratedNode()
	pageRouter.Add(PAGE_INTEGRATED_NODE, pageIntegratedNode)

	pageRemoteNode := NewPageRemoteNode()
	pageRouter.Add(PAGE_REMOTE_NODE, pageRemoteNode)

	header := prefabs.NewHeader(pageRouter)

	page := &Page{
		pageRouter:         pageRouter,
		pageSelectNode:     pageSelectNode,
		pageAddNodeForm:    pageAddNodeForm,
		pageEditNodeForm:   pageEditNodeForm,
		pageRemoteNode:     pageRemoteNode,
		pageIntegratedNode: pageIntegratedNode,
		header:             header,

		pageSectionAnimation: pages.NewPageSectionAnimation(),
	}
	page_instance = page
	return page
}

func (p *Page) IsActive() bool {
	return p.isActive
}

func (p *Page) Enter() {
	bottom_bar.Instance.SetButtonActive(bottom_bar.BUTTON_NODE)
	p.isActive = p.pageSectionAnimation.Enter()

	currentNode := node_manager.CurrentNode

	p.header.ResetHistory()
	p.header.AddHistory(PAGE_SELECT_NODE)
	if currentNode != nil {
		if currentNode.Integrated {
			p.header.AddHistory(PAGE_INTEGRATED_NODE)
			//p.pageIntegratedNode.animationLeave.Reset()
			p.pageRouter.SetCurrent(PAGE_INTEGRATED_NODE)
		} else {
			p.header.AddHistory(PAGE_REMOTE_NODE)
			//p.pageRemoteNode.animationLeave.Reset()
			p.pageRouter.SetCurrent(PAGE_REMOTE_NODE)
		}
	} else {
		p.pageRouter.SetCurrent(PAGE_SELECT_NODE)
	}
}

func (p *Page) Leave() {
	//p.pageRemoteNode.animationLeave.Reset()
	//p.pageIntegratedNode.animationLeave.Reset()
	p.pageSectionAnimation.Leave()
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
