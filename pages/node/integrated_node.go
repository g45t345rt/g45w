package page_node

import (
	"fmt"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/node"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/utils"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageIntegratedNode struct {
	isActive       bool
	animationEnter *animation.Animation
	animationLeave *animation.Animation

	hubIcon    *widget.Icon
	nodeSize   *node.NodeSize
	nodeStatus *node.NodeStatus
}

var _ router.Container = &PageIntegratedNode{}

func NewPageIntegratedNode() *PageIntegratedNode {
	hubIcon, _ := widget.NewIcon(icons.HardwareDeviceHub)

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .5, ease.OutCubic),
	))

	return &PageIntegratedNode{
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		hubIcon:        hubIcon,
		nodeSize:       node.NewNodeSize(10 * time.Second),
		nodeStatus:     node.NewNodeStatus(1 * time.Second),
	}
}

func (p *PageIntegratedNode) IsActive() bool {
	return p.isActive
}

func (p *PageIntegratedNode) Enter() {
	p.isActive = true

	page_instance.header.SetTitle("Select Node")
	p.animationLeave.Reset()
	p.animationEnter.Start()
}

func (p *PageIntegratedNode) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *PageIntegratedNode) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	p.nodeStatus.Active()
	p.nodeSize.Active()

	return layout.Inset{
		Top: unit.Dp(30), Bottom: unit.Dp(30),
		Left: unit.Dp(30), Right: unit.Dp(30),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(30)
						gtx.Constraints.Min.Y = gtx.Dp(30)

						return p.hubIcon.Layout(gtx, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Label(th, unit.Sp(24), "NODE")
						label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
						return label.Layout(gtx)
					}),
				)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), "Node Height / Network Height")
				label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 100}
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				status := fmt.Sprintf("%d / %d", p.nodeStatus.Height, p.nodeStatus.BestHeight)
				label := material.Label(th, unit.Sp(22), status)
				label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
				return label.Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), "Peers")
				label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 100}
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				status := fmt.Sprintf("%d", p.nodeStatus.PeerCount)
				label := material.Label(th, unit.Sp(22), status)
				label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
				return label.Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), "Network Hashrate")
				label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 100}
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				status := utils.FormatHashRate(p.nodeStatus.NetworkHashRate)
				label := material.Label(th, unit.Sp(22), status)
				label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
				return label.Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), "TXp / Time Offset")
				label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 100}
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				status := fmt.Sprintf("%d:%d / %s | %s | %s",
					p.nodeStatus.MemCount,
					p.nodeStatus.RegCount,
					p.nodeStatus.TimeOffset.String(),
					p.nodeStatus.TimeOffsetNTP.String(),
					p.nodeStatus.TimeOffsetP2P.String(),
				)

				label := material.Label(th, unit.Sp(22), status)
				label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
				return label.Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), "Space Used")
				label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 100}
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				value := utils.FormatBytes(p.nodeSize.Size)
				label := material.Label(th, unit.Sp(22), value)
				label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
				return label.Layout(gtx)
			}),
		)
	})
}
