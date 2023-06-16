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

	currentNode := node_manager.Instance.NodeState.Current
	status := "unassigned node"
	if currentNode != "" {
		if currentNode == node_manager.INTEGRATED_NODE_ID {
			n.integratedNodeStatus.Active()

			height := n.integratedNodeStatus.Height
			bestHeight := n.integratedNodeStatus.BestHeight
			out := n.integratedNodeStatus.PeerOutCount
			status = fmt.Sprintf("%d / %d - %dP (%s)", height, bestHeight, out, "Integrated")
		} else {
			n.remoteNodeInfo.Active()
			nodeConn := node_manager.Instance.NodeState.Nodes[currentNode]
			height := n.remoteNodeInfo.Result.Height
			out := n.remoteNodeInfo.Result.Outgoing_connections_count
			status = fmt.Sprintf("%d - %dP (%s)", height, out, nodeConn.Name)
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
					gtx.Constraints.Max.X = gtx.Dp(12)
					gtx.Constraints.Max.Y = gtx.Dp(12)
					paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 0, B: 0, A: 255},
						clip.Ellipse{
							Max: gtx.Constraints.Max,
						}.Op(gtx.Ops))

					return layout.Dimensions{Size: gtx.Constraints.Max}
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					label := material.Label(th, unit.Sp(16), status)
					label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
					return label.Layout(gtx)
				}),
			)
		})
	})
}
