package page_node

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
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
	router        *router.Router

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
		Rounded:         unit.Dp(5),
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

	router := router.NewRouter()
	pageSelectNode := NewPageSelectNode()
	router.Add(PAGE_SELECT_NODE, pageSelectNode)

	pageAddNodeForm := NewPageAddNodeForm()
	router.Add(PAGE_ADD_NODE_FORM, pageAddNodeForm)

	pageEditNodeForm := NewPageEditNodeForm()
	router.Add(PAGE_EDIT_NODE_FORM, pageEditNodeForm)

	pageIntegratedNode := NewPageIntegratedNode()
	router.Add(PAGE_INTEGRATED_NODE, pageIntegratedNode)

	pageRemoteNode := NewPageRemoteNode()
	router.Add(PAGE_REMOTE_NODE, pageRemoteNode)

	th := app_instance.Theme
	labelHeaderStyle := material.Label(th, unit.Sp(22), "")
	labelHeaderStyle.Font.Weight = font.Bold
	header := prefabs.NewHeader(labelHeaderStyle, router, nil)

	page := &Page{
		buttonSetNode:    buttonSetNode,
		firstEnter:       true,
		router:           router,
		pageSelectNode:   pageSelectNode,
		pageAddNodeForm:  pageAddNodeForm,
		pageEditNodeForm: pageEditNodeForm,
		pageRemoteNode:   pageRemoteNode,
		header:           header,

		animationEnter: animationEnter,
		animationLeave: animationLeave,
	}
	page_instance = page
	router.SetPrimary(PAGE_SELECT_NODE)

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

	currentNode := node_manager.Instance.NodeState.Current
	if currentNode != "" {
		if currentNode == node_manager.INTEGRATED_NODE_ID {
			p.router.SetCurrent(PAGE_INTEGRATED_NODE)
		} else {
			p.router.SetCurrent(PAGE_REMOTE_NODE)
		}
	}
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

	if p.pageSelectNode.buttonAddNode.Clickable.Clicked() {
		p.router.SetCurrent(PAGE_ADD_NODE_FORM)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top: unit.Dp(30), Bottom: unit.Dp(20),
				Left: unit.Dp(30), Right: unit.Dp(30),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return p.header.Layout(gtx, th, nil)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return p.router.Layout(gtx, th)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return bottom_bar.Instance.Layout(gtx, th)
		}),
	)
}
