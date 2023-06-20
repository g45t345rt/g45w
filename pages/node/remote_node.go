package page_node

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/node_manager"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/utils"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type PageRemoteNode struct {
	isActive          bool
	useAnimationEnter bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	nodeInfo *RemoteNodeInfo
}

var _ router.Page = &PageRemoteNode{}

func NewPageRemoteNode() *PageRemoteNode {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .5, ease.OutCubic),
	))

	return &PageRemoteNode{
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		nodeInfo:       NewRemoteNodeInfo(3 * time.Second),
	}
}

func (p *PageRemoteNode) IsActive() bool {
	return p.isActive
}

func (p *PageRemoteNode) Enter() {
	p.isActive = true

	page_instance.header.SetTitle("Remote Node")

	if p.useAnimationEnter {
		p.animationLeave.Reset()
		p.animationEnter.Start()
	}
}

func (p *PageRemoteNode) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *PageRemoteNode) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	currentNode := node_manager.Instance.NodeState.Current
	nodeInfo := node_manager.Instance.NodeState.Nodes[currentNode]

	p.nodeInfo.Active()

	return layout.Inset{
		Top: unit.Dp(0), Bottom: unit.Dp(30),
		Left: unit.Dp(30), Right: unit.Dp(30),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				r := op.Record(gtx.Ops)
				dims := layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Label(th, unit.Sp(22), nodeInfo.Name)
							label.Color = color.NRGBA{A: 255}
							return label.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Label(th, unit.Sp(16), nodeInfo.Endpoint)
							label.Color = color.NRGBA{A: 150}
							return label.Layout(gtx)
						}),
					)
				})
				c := r.Stop()

				paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255},
					clip.UniformRRect(
						image.Rectangle{Max: dims.Size},
						gtx.Dp(15),
					).Op(gtx.Ops))

				c.Add(gtx.Ops)
				return dims
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), "Node Height")
				label.Color = color.NRGBA{A: 150}
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				status := fmt.Sprintf("%d", p.nodeInfo.Result.Height)
				label := material.Label(th, unit.Sp(22), status)
				label.Color = color.NRGBA{A: 255}
				return label.Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), "Peers (In/Out)")
				label.Color = color.NRGBA{A: 150}
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				inc := p.nodeInfo.Result.Incoming_connections_count
				out := p.nodeInfo.Result.Outgoing_connections_count
				status := fmt.Sprintf("%d / %d", inc, out)
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
				diff := p.nodeInfo.Result.Difficulty
				status := utils.FormatHashRate(diff)
				label := material.Label(th, unit.Sp(22), status)
				label.Color = color.NRGBA{A: 255}
				return label.Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), "Version")
				label.Color = color.NRGBA{A: 150}
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				version := p.nodeInfo.Result.Version
				label := material.Label(th, unit.Sp(22), version)
				label.Color = color.NRGBA{A: 255}
				return label.Layout(gtx)
			}),
		)
	})
}

type RemoteNodeInfo struct {
	Result rpc.GetInfo_Result

	isActive bool
}

func NewRemoteNodeInfo(d time.Duration) *RemoteNodeInfo {
	nodeInfo := &RemoteNodeInfo{isActive: false}
	ticker := time.NewTicker(d)

	window := app_instance.Window
	go func() {
		for range ticker.C {
			if nodeInfo.isActive {
				nodeInfo.update()
				window.Invalidate()
				nodeInfo.isActive = false
			}
		}
	}()

	nodeInfo.update()
	return nodeInfo
}

func (n *RemoteNodeInfo) Active() {
	n.isActive = true
}

func (n *RemoteNodeInfo) update() {
	if walletapi.RPC_Client.RPC == nil {
		return
	}

	walletapi.RPC_Client.RPC.CallResult(context.Background(), "DERO.GetInfo", nil, &n.Result)
}
