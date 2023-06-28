package page_node

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/containers/bottom_bar"
	"github.com/g45t345rt/g45w/node_manager"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Page struct {
	isActive   bool
	firstEnter bool

	buttonSetNode *components.Button
	pageRouter    *router.Router

	pageSelectNode   *PageSelectNode
	pageAddNodeForm  *PageAddNodeForm
	pageEditNodeForm *PageEditNodeForm
	pageRemoteNode   *PageRemoteNode
	header           *prefabs.Header

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

	setIcon, _ := widget.NewIcon(icons.ActionSettings)
	buttonSetNode := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Text:            "SELECT NODE",
		Icon:            setIcon,
		TextColor:       color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		BackgroundColor: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonSetNode.Label.Alignment = text.Middle
	buttonSetNode.Style.Font.Weight = font.Bold

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

	th := app_instance.Theme
	labelHeaderStyle := material.Label(th, unit.Sp(22), "")
	labelHeaderStyle.Font.Weight = font.Bold
	header := prefabs.NewHeader(labelHeaderStyle, pageRouter)

	page := &Page{
		buttonSetNode:    buttonSetNode,
		firstEnter:       true,
		pageRouter:       pageRouter,
		pageSelectNode:   pageSelectNode,
		pageAddNodeForm:  pageAddNodeForm,
		pageEditNodeForm: pageEditNodeForm,
		pageRemoteNode:   pageRemoteNode,
		header:           header,

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
	if currentNode != "" {
		if p.pageRouter.Current == nil {
			p.header.AddHistory(PAGE_SELECT_NODE)
		}

		if currentNode == node_manager.INTEGRATED_NODE_ID {
			p.pageRouter.SetCurrent(PAGE_INTEGRATED_NODE)
		} else {
			p.pageRouter.SetCurrent(PAGE_REMOTE_NODE)
		}
	} else {
		p.pageRouter.SetCurrent(PAGE_SELECT_NODE)
	}
}

func (p *Page) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.pageSelectNode.buttonAddNode.Clickable.Clicked() {
		p.pageRouter.SetCurrent(PAGE_ADD_NODE_FORM)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top: unit.Dp(30), Bottom: unit.Dp(20),
				Left: unit.Dp(30), Right: unit.Dp(30),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return p.header.Layout(gtx, th)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return p.pageRouter.Layout(gtx, th)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return bottom_bar.Instance.Layout(gtx, th)
		}),
	)
}
