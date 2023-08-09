package page_node

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/containers/bottom_bar"
	"github.com/g45t345rt/g45w/node_manager"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
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

	animationEnter *animation.Animation
	animationLeave *animation.Animation
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
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .5, ease.OutCubic),
	))

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

		animationEnter: animationEnter,
		animationLeave: animationLeave,
	}
	page_instance = page
	return page
}

func (p *Page) IsActive() bool {
	return p.isActive
}

func (p *Page) Enter() {
	bottom_bar.Instance.SetButtonActive(bottom_bar.BUTTON_NODE)
	p.isActive = true
	p.animationEnter.Start()
	p.animationLeave.Reset()

	currentNode := node_manager.CurrentNode

	p.header.ResetHistory()
	p.header.AddHistory(PAGE_SELECT_NODE)
	if currentNode != nil {
		if currentNode.Integrated {
			p.header.AddHistory(PAGE_INTEGRATED_NODE)
			p.pageIntegratedNode.animationLeave.Reset()
			p.pageRouter.SetCurrent(PAGE_INTEGRATED_NODE)
		} else {
			p.header.AddHistory(PAGE_REMOTE_NODE)
			p.pageRemoteNode.animationLeave.Reset()
			p.pageRouter.SetCurrent(PAGE_REMOTE_NODE)
		}
	} else {
		p.pageRouter.SetCurrent(PAGE_SELECT_NODE)
	}
}

func (p *Page) Leave() {
	p.pageRemoteNode.animationLeave.Reset()
	p.pageIntegratedNode.animationLeave.Reset()

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

			startColor := theme.Current.BgGradientStartColor
			endColor := theme.Current.BgGradientEndColor
			defer utils.PaintLinearGradient(gtx, startColor, endColor).Pop()

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top: unit.Dp(30), Bottom: unit.Dp(20),
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
