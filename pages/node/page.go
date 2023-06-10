package page_node

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/pages"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/components"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Page struct {
	isActive   bool
	firstEnter bool

	buttonSetNode *components.Button
	childRouter   *router.Router

	pageSelectNode   *PageSelectNode
	pageAddNodedForm *PageAddNodeForm
	header           *prefabs.Header
}

var _ router.Container = &Page{}

type PageInstance struct {
	router *router.Router
	header *prefabs.Header
}

var page_instance *PageInstance

func NewPage() *Page {
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

	childRouter := router.NewRouter()
	pageSelectNode := NewPageSelectNode()
	childRouter.Add("selectNode", pageSelectNode)

	pageAddNodeForm := NewPageAddNodeForm()
	childRouter.Add("addNodeForm", pageAddNodeForm)

	pageIntegratedNode := NewPageIntegratedNode()
	childRouter.Add("integratedNode", pageIntegratedNode)

	th := app_instance.Current.Theme
	labelHeaderStyle := material.Label(th, unit.Sp(22), "")
	labelHeaderStyle.Font.Weight = font.Bold
	header := prefabs.NewHeader(labelHeaderStyle, childRouter, nil)

	page_instance = &PageInstance{
		router: childRouter,
		header: header,
	}

	childRouter.SetPrimary("selectNode")

	return &Page{
		buttonSetNode:    buttonSetNode,
		firstEnter:       true,
		childRouter:      childRouter,
		pageSelectNode:   pageSelectNode,
		pageAddNodedForm: pageAddNodeForm,
		header:           header,
	}
}

func (p *Page) IsActive() bool {
	return p.isActive
}

func (p *Page) Enter() {
	pages.BottomBarInstance.SetButtonActive("node")
	p.isActive = true
}

func (p *Page) Leave() {
	p.isActive = false
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

	if p.pageSelectNode.buttonAddNode.Clickable.Clicked() {
		p.childRouter.SetCurrent("addNodeForm")
	}

	if p.pageSelectNode.buttonSetIntegratedNode.Clickable.Clicked() {
		p.childRouter.SetCurrent("integratedNode")
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
			return p.childRouter.Layout(gtx, th)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return pages.BottomBarInstance.Layout(gtx, th)
		}),
	)
}
