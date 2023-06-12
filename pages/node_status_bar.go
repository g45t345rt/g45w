package pages

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
)

type NodeStatusBar struct {
	clickable  *widget.Clickable
	nodeStatus *integrated_node.NodeStatus
}

var NodeStatusBarInstance *NodeStatusBar

func LoadNodeStatusBarInstance() *NodeStatusBar {
	nodeStatusBar := &NodeStatusBar{
		clickable:  new(widget.Clickable),
		nodeStatus: integrated_node.NewNodeStatus(1 * time.Second),
	}
	NodeStatusBarInstance = nodeStatusBar
	return nodeStatusBar
}

func (n *NodeStatusBar) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	paint.FillShape(gtx.Ops, color.NRGBA{A: 255}, clip.Rect{
		Max: gtx.Constraints.Max,
	}.Op())

	//paint.ColorOp{Color: color.NRGBA{A: 255}}.Add(gtx.Ops)
	//paint.PaintOp{}.Add(gtx.Ops)

	n.nodeStatus.Active()

	if n.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
	}

	if n.clickable.Clicked() {
		app_instance.Router.SetCurrent("page_node")
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
					status := fmt.Sprintf("%d / %d - %dP", n.nodeStatus.Height, n.nodeStatus.BestHeight, n.nodeStatus.PeerCount)
					label := material.Label(th, unit.Sp(16), status)
					label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
					return label.Layout(gtx)
				}),
			)
		})
	})
}
