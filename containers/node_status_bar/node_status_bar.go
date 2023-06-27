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
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/integrated_node"
	"github.com/g45t345rt/g45w/node_manager"
	page_node "github.com/g45t345rt/g45w/pages/node"
	"github.com/g45t345rt/g45w/wallet_manager"
)

type NodeStatusBar struct {
	clickable            *widget.Clickable
	integratedNodeStatus *integrated_node.NodeStatus
	remoteNodeInfo       *page_node.RemoteNodeInfo
}

var Instance *NodeStatusBar

func LoadInstance() *NodeStatusBar {
	nodeStatusBar := &NodeStatusBar{
		clickable:            new(widget.Clickable),
		integratedNodeStatus: integrated_node.NewNodeStatus(1 * time.Second),
		remoteNodeInfo:       page_node.NewRemoteNodeInfo(3 * time.Second),
	}
	Instance = nodeStatusBar
	return nodeStatusBar
}

func (n *NodeStatusBar) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	paint.FillShape(gtx.Ops, color.NRGBA{A: 255}, clip.Rect{
		Max: gtx.Constraints.Max,
	}.Op())

	//paint.ColorOp{Color: color.NRGBA{A: 255}}.Add(gtx.Ops)
	//paint.PaintOp{}.Add(gtx.Ops)
	if wallet_manager.OpenedWallet == nil {
		return layout.Dimensions{}
	}

	wallet := wallet_manager.OpenedWallet.Memory
	currentNode := node_manager.CurrentNode
	status := "unassigned node"
	statusDotColor := color.NRGBA{R: 255, G: 0, B: 0, A: 255}

	if currentNode != "" {
		if currentNode == node_manager.INTEGRATED_NODE_ID {
			n.integratedNodeStatus.Active()

			//height := n.integratedNodeStatus.Height
			//bestHeight := n.integratedNodeStatus.BestHeight
			walletHeight := wallet.Get_Height()
			daemonHeight := wallet.Get_Daemon_Height()
			out := n.integratedNodeStatus.PeerOutCount

			if walletHeight < daemonHeight {
				statusDotColor = color.NRGBA{R: 255, G: 255, B: 0, A: 255}
			} else {
				statusDotColor = color.NRGBA{R: 0, G: 255, B: 0, A: 255}
			}

			status = fmt.Sprintf("%d / %d - %dP (%s)", walletHeight, daemonHeight, out, "Integrated")
		} else {
			n.remoteNodeInfo.Active()

			nodeConn := node_manager.Nodes[currentNode]
			walletHeight := wallet.Get_Height()
			daemonHeight := wallet.Get_Daemon_Height()
			out := n.remoteNodeInfo.Result.Outgoing_connections_count

			if walletHeight < daemonHeight {
				statusDotColor = color.NRGBA{R: 255, G: 255, B: 0, A: 255}
			} else {
				statusDotColor = color.NRGBA{R: 0, G: 255, B: 0, A: 255}
			}

			status = fmt.Sprintf("%d / %d - %dP (%s)", walletHeight, daemonHeight, out, nodeConn.Name)
		}
	}

	if n.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
	}

	if n.clickable.Clicked() {
		app_instance.Router.SetCurrent(app_instance.PAGE_NODE)
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	return n.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(15), Bottom: unit.Dp(15),
			Left: unit.Dp(10), Right: unit.Dp(10),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return StatusDot{
						Color: statusDotColor,
					}.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), status)
					lbl.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
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
