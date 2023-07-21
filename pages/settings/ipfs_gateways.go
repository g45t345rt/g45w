package page_settings

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageIPFSGateways struct {
	isActive       bool
	list           *widget.List
	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonAdd *components.Button
}

var _ router.Page = &PageIPFSGateways{}

func NewPageIPFSGateways() *PageIPFSGateways {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	addIcon, _ := widget.NewIcon(icons.ContentAdd)
	buttonAdd := components.NewButton(components.ButtonStyle{
		Icon:      addIcon,
		TextColor: color.NRGBA{A: 255},
	})

	return &PageIPFSGateways{
		list:           list,
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonAdd: buttonAdd,
	}
}

func (p *PageIPFSGateways) IsActive() bool {
	return p.isActive
}

func (p *PageIPFSGateways) Enter() {
	p.isActive = true
	page_instance.header.SetTitle(lang.Translate("IPFS Gateways"))
	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(16), lang.Translate("Interplanetary File System"))
		return lbl.Layout(gtx)
	}
	page_instance.header.ButtonRight = p.buttonAdd

	if !page_instance.header.IsHistory(PAGE_IPFS_GATEWAYS) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageIPFSGateways) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *PageIPFSGateways) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	{
		state := p.animationEnter.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	{
		state := p.animationLeave.Update(gtx)
		if state.Finished {
			p.isActive = false
			op.InvalidateOp{}.Add(gtx.Ops)
		}

		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	widgets := []layout.Widget{}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}
