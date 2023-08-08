package page_node

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/node_manager"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/utils"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageRemoteNode struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonReconnect *components.Button
	nodeInfo        *RemoteNodeInfo
	connecting      bool
}

var _ router.Page = &PageRemoteNode{}

func NewPageRemoteNode() *PageRemoteNode {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .5, ease.OutCubic),
	))

	refreshIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	buttonReconnect := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            refreshIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonReconnect.Label.Alignment = text.Middle
	buttonReconnect.Style.Font.Weight = font.Bold

	nodeInfo := NewRemoteNodeInfo(3 * time.Second)

	return &PageRemoteNode{
		animationEnter:  animationEnter,
		animationLeave:  animationLeave,
		nodeInfo:        nodeInfo,
		buttonReconnect: buttonReconnect,
	}
}

func (p *PageRemoteNode) IsActive() bool {
	return p.isActive
}

func (p *PageRemoteNode) Enter() {
	p.isActive = true
	page_instance.header.Title = func() string { return lang.Translate("Remote Node") }

	if !page_instance.header.IsHistory(PAGE_REMOTE_NODE) {
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

	currentNode := node_manager.CurrentNode
	if currentNode == nil {
		return layout.Dimensions{}
	}

	p.nodeInfo.Active()

	if p.buttonReconnect.Clicked() {
		p.reconnect()
	}

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
							label := material.Label(th, unit.Sp(22), currentNode.Name)
							label.Color = color.NRGBA{A: 255}
							return label.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Label(th, unit.Sp(16), currentNode.Endpoint)
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
				if p.nodeInfo.Err != nil {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(18), lang.Translate("Error"))
							lbl.Font.Weight = font.Bold
							return lbl.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(16), p.nodeInfo.Err.Error())
							return lbl.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							p.buttonReconnect.Text = lang.Translate("Reconnect")
							return p.buttonReconnect.Layout(gtx, th)
						}),
					)
				}

				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Label(th, unit.Sp(18), lang.Translate("Node Height"))
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
						label := material.Label(th, unit.Sp(18), lang.Translate("Peers (In/Out)"))
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
						label := material.Label(th, unit.Sp(18), lang.Translate("Network Hashrate"))
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
						label := material.Label(th, unit.Sp(18), lang.Translate("Version"))
						label.Color = color.NRGBA{A: 150}
						return label.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						version := p.nodeInfo.Result.Version
						label := material.Label(th, unit.Sp(16), version)
						label.Color = color.NRGBA{A: 255}
						return label.Layout(gtx)
					}),
				)
			}),
		)
	})
}

func (p *PageRemoteNode) reconnect() {
	if p.connecting {
		return
	}

	p.connecting = true
	go func() {
		currentNode := node_manager.CurrentNode

		notification_modals.InfoInstance.SetText(lang.Translate("Connecting..."), currentNode.Endpoint)
		notification_modals.InfoInstance.SetVisible(true, 0)
		err := node_manager.Connect(*currentNode, true)
		p.connecting = false
		notification_modals.InfoInstance.SetVisible(false, 0)

		if err != nil {
			notification_modals.InfoInstance.SetVisible(false, 0)
			notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {
			p.nodeInfo.Update()
			app_instance.Window.Invalidate()
			notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("Remote node reconnected."))
			notification_modals.SuccessInstance.SetVisible(true, 0)
		}
	}()
}

type RemoteNodeInfo struct {
	Result rpc.GetInfo_Result
	Err    error

	isActive bool
}

func NewRemoteNodeInfo(d time.Duration) *RemoteNodeInfo {
	nodeInfo := &RemoteNodeInfo{isActive: false}
	ticker := time.NewTicker(d)

	window := app_instance.Window
	go func() {
		for range ticker.C {
			if nodeInfo.isActive {
				nodeInfo.Update()
				window.Invalidate()
				nodeInfo.isActive = false
			}
		}
	}()

	nodeInfo.Update()
	return nodeInfo
}

func (n *RemoteNodeInfo) Active() {
	n.isActive = true
}

func (n *RemoteNodeInfo) Update() {
	if walletapi.RPC_Client.RPC == nil {
		return
	}

	n.Err = walletapi.RPC_Client.RPC.CallResult(context.Background(), "DERO.GetInfo", nil, &n.Result)
}
