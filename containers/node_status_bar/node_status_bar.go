package node_status_bar

import (
	"fmt"
	"image/color"
	"time"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/integrated_node"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/node_manager"
	"github.com/g45t345rt/g45w/pages"
	page_node "github.com/g45t345rt/g45w/pages/node"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type NodeStatusBar struct {
	clickable            *widget.Clickable
	IntegratedNodeStatus *integrated_node.NodeStatus
	RemoteNodeInfo       *page_node.RemoteNodeInfo
	pulseAnimation       *animation.Animation
}

var Instance *NodeStatusBar

func LoadInstance() *NodeStatusBar {
	nodeStatusBar := &NodeStatusBar{
		clickable:            new(widget.Clickable),
		IntegratedNodeStatus: integrated_node.NewNodeStatus(1 * time.Second),
		RemoteNodeInfo:       page_node.NewRemoteNodeInfo(3 * time.Second),
		pulseAnimation: animation.NewAnimation(false, gween.NewSequence(
			gween.New(1, .7, .5, ease.OutBounce),
			gween.New(.7, 1, .5, ease.InBounce),
		)),
	}
	Instance = nodeStatusBar
	return nodeStatusBar
}

func (n *NodeStatusBar) Update() {
	currentNode := node_manager.CurrentNode
	if currentNode != nil {
		if currentNode.Integrated {
			n.IntegratedNodeStatus.Update()
		} else {
			n.RemoteNodeInfo.Update()
		}
	}
}

func (n *NodeStatusBar) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	bgColor := theme.Current.NodeStatusBgColor
	paint.FillShape(gtx.Ops, bgColor, clip.Rect{
		Max: gtx.Constraints.Max,
	}.Op())

	if wallet_manager.OpenedWallet == nil {
		return layout.Dimensions{}
	}

	wallet := wallet_manager.OpenedWallet
	currentNode := node_manager.CurrentNode
	status := lang.Translate("Unassigned Node")
	statusDotColor := theme.Current.NodeStatusDotRedColor // color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	nodeName := ""

	if currentNode != nil {
		if currentNode.Integrated {
			n.IntegratedNodeStatus.Active()

			//height := n.integratedNodeStatus.Height
			//bestHeight := n.integratedNodeStatus.BestHeight
			walletHeight := wallet.Memory.Get_Height()
			daemonHeight := wallet.Memory.Get_Daemon_Height()
			//out := n.IntegratedNodeStatus.PeerOutCount

			if walletHeight < daemonHeight {
				statusDotColor = theme.Current.NodeStatusDotYellowColor
			} else {
				statusDotColor = theme.Current.NodeStatusDotGreenColor
			}

			nodeName = lang.Translate("Integrated Node")
			status = fmt.Sprintf("%d / %d", walletHeight, daemonHeight)
		} else {
			n.RemoteNodeInfo.Active()
			walletHeight := wallet.Memory.Get_Height()
			daemonHeight := wallet.Memory.Get_Daemon_Height()
			//out := n.RemoteNodeInfo.Result.Outgoing_connections_count
			nodeName = currentNode.Name

			if n.RemoteNodeInfo.Err == nil {
				if walletHeight < daemonHeight {
					statusDotColor = theme.Current.NodeStatusDotYellowColor //color.NRGBA{R: 255, G: 255, B: 0, A: 255}
				} else {
					statusDotColor = theme.Current.NodeStatusDotGreenColor // color.NRGBA{R: 0, G: 255, B: 0, A: 255}
				}

				status = fmt.Sprintf("%d / %d", walletHeight, daemonHeight)
			} else {
				status = lang.Translate("Disconnected")
			}
		}
	}

	if n.clickable.Clicked(gtx) {
		app_instance.Router.SetCurrent(pages.PAGE_NODE)
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	return n.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if n.clickable.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
		}

		return layout.Inset{
			Top: unit.Dp(15), Bottom: unit.Dp(15),
			Left: unit.Dp(10), Right: unit.Dp(10),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// pulse if yellow
					if statusDotColor == theme.Current.NodeStatusDotYellowColor {
						n.pulseAnimation.Start()
					} else {
						n.pulseAnimation.Reset()
					}

					r := op.Record(gtx.Ops)
					dims := StatusDot{
						Color: statusDotColor,
					}.Layout(gtx)
					c := r.Stop()
					gtx.Constraints.Min = dims.Size

					state := n.pulseAnimation.Update(gtx)
					if state.Active {
						defer animation.TransformScaleCenter(gtx, state.Value).Push(gtx.Ops).Pop()
					}

					c.Add(gtx.Ops)
					return dims
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), status)
					lbl.Color = theme.Current.NodeStatusHeightTextColor
					return lbl.Layout(gtx)
				}),
				//layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Spacer{Width: unit.Dp(1)}.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), nodeName)
					lbl.Color = theme.Current.NodeStatusNodeTextColor
					return lbl.Layout(gtx)
				}),
			)
		})
	})
}

type StatusDot struct {
	Color color.NRGBA
}

func (s StatusDot) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.X = gtx.Dp(12)
	gtx.Constraints.Max.Y = gtx.Dp(12)
	paint.FillShape(gtx.Ops, s.Color,
		clip.Ellipse{
			Max: gtx.Constraints.Max,
		}.Op(gtx.Ops))

	return layout.Dimensions{Size: gtx.Constraints.Max}
}
