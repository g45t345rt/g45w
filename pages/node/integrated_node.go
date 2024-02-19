package page_node

import (
	"fmt"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/integrated_node"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
)

type PageIntegratedNode struct {
	isActive            bool
	headerPageAnimation *prefabs.PageHeaderAnimation

	nodeStatus     integrated_node.NodeStatus
	nodeStatusLoop *utils.ForceActiveLoop
	nodeSize       int64
	nodeSizeLoop   *utils.ForceActiveLoop
}

var _ router.Page = &PageIntegratedNode{}

func NewPageIntegratedNode() *PageIntegratedNode {
	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_INTEGRATED_NODE)

	page := &PageIntegratedNode{
		headerPageAnimation: headerPageAnimation,
	}

	page.nodeStatusLoop = utils.NewForceActiveLoop(1*time.Second, func() {
		page.nodeStatus = integrated_node.GetStatus()
	})

	page.nodeSizeLoop = utils.NewForceActiveLoop(10*time.Second, func() {
		nodeDir := settings.IntegratedNodeDir
		size, _ := utils.GetFolderSize(nodeDir)
		page.nodeSize = size
	})

	return page
}

func (p *PageIntegratedNode) IsActive() bool {
	return p.isActive
}

func (p *PageIntegratedNode) Enter() {
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)
	page_instance.header.Title = func() string { return lang.Translate("Integrated Node") }
}

func (p *PageIntegratedNode) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)
}

func (p *PageIntegratedNode) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	p.nodeStatusLoop.SetActive()
	p.nodeSizeLoop.SetActive()

	return layout.Inset{
		Top: unit.Dp(0), Bottom: unit.Dp(30),
		Left: theme.PagePadding, Right: theme.PagePadding,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), lang.Translate("Node Height / Network Height"))
				label.Color = theme.Current.TextMuteColor
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				status := fmt.Sprintf("%d / %d", p.nodeStatus.Height, p.nodeStatus.BestHeight)
				label := material.Label(th, unit.Sp(22), status)
				return label.Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), lang.Translate("Peers (In/Out)"))
				label.Color = theme.Current.TextMuteColor
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				status := fmt.Sprintf("%d / %d", p.nodeStatus.PeerInCount, p.nodeStatus.PeerOutCount)
				label := material.Label(th, unit.Sp(22), status)
				return label.Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), lang.Translate("Network Hashrate"))
				label.Color = theme.Current.TextMuteColor
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				status := utils.FormatHashRate(p.nodeStatus.NetworkHashRate)
				label := material.Label(th, unit.Sp(22), status)
				return label.Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), lang.Translate("TXp / Time Offset"))
				label.Color = theme.Current.TextMuteColor
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
				return label.Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), lang.Translate("Space Used"))
				label.Color = theme.Current.TextMuteColor
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				value := utils.FormatBytes(p.nodeSize)
				label := material.Label(th, unit.Sp(22), value)
				return label.Layout(gtx)
			}),
		)
	})
}
