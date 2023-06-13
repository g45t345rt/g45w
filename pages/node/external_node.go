package page_node

import (
	"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/utils"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type PageExternalNode struct {
	isActive       bool
	animationEnter *animation.Animation
	animationLeave *animation.Animation
}

var _ router.Page = &PageExternalNode{}

func NewPageExternalNode() *PageExternalNode {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .5, ease.OutCubic),
	))

	return &PageExternalNode{
		animationEnter: animationEnter,
		animationLeave: animationLeave,
	}
}

func (p *PageExternalNode) IsActive() bool {
	return p.isActive
}

func (p *PageExternalNode) Enter() {
	p.isActive = true

	page_instance.header.SetTitle("External Node")
	p.animationLeave.Reset()
	p.animationEnter.Start()
}

func (p *PageExternalNode) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *PageExternalNode) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	return layout.Inset{
		Top: unit.Dp(0), Bottom: unit.Dp(30),
		Left: unit.Dp(30), Right: unit.Dp(30),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), "Node Height / Network Height")
				label.Color = color.NRGBA{A: 150}
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				status := fmt.Sprintf("%d / %d", 0, 0)
				label := material.Label(th, unit.Sp(22), status)
				label.Color = color.NRGBA{A: 255}
				return label.Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), "Peers")
				label.Color = color.NRGBA{A: 150}
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				status := fmt.Sprintf("%d", 0)
				label := material.Label(th, unit.Sp(22), status)
				label.Color = color.NRGBA{A: 255}
				return label.Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), "Network Hashrate")
				label.Color = color.NRGBA{A: 150}
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				status := utils.FormatHashRate(0)
				label := material.Label(th, unit.Sp(22), status)
				label.Color = color.NRGBA{A: 255}
				return label.Layout(gtx)
			}),
		)
	})
}
