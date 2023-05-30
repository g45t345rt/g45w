package page_node

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/p2p"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/g45t345rt/g45w/utils"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Page struct {
	isActive   bool
	firstEnter bool

	buttonShowCommand *components.Button
	hubIcon           *widget.Icon
	nodeSize          *NodeSize
}

var _ router.Container = &Page{}

func NewPage() *Page {
	hubIcon, _ := widget.NewIcon(icons.HardwareDeviceHub)

	cmdIcon, _ := widget.NewIcon(icons.ActionTouchApp)
	buttonShowCommand := components.NewButton(components.ButtonStyle{
		Rounded:         unit.Dp(5),
		Text:            "COMMANDS",
		Icon:            cmdIcon,
		TextColor:       color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		BackgroundColor: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonShowCommand.Label.Alignment = text.Middle
	buttonShowCommand.Style.Font.Weight = font.Bold

	return &Page{
		buttonShowCommand: buttonShowCommand,
		hubIcon:           hubIcon,
		firstEnter:        true,
		nodeSize:          NewNodeSize(),
	}
}

func (p *Page) IsActive() bool {
	return p.isActive
}

func (p *Page) Enter() {
	app_instance.Current.BottomBar.SetActive("node")
	p.isActive = true
	p.nodeSize.active = true
}

func (p *Page) Leave() {
	p.isActive = false
	p.nodeSize.active = false
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	paint.ColorOp{Color: color.NRGBA{A: 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	chain := app_instance.Current.Chain
	our_height := chain.Get_Height()
	best_height, _ := p2p.Best_Peer_Height()
	//topo_height := chain.Load_TOPO_HEIGHT()

	mempool_tx_count := len(chain.Mempool.Mempool_List_TX())
	regpool_tx_count := len(chain.Regpool.Regpool_List_TX())

	//p2p.PeerList_Print()
	peer_count := p2p.Peer_Count()
	//inc, out := p2p.Peer_Direction_Count()
	network_hash_rate := chain.Get_Network_HashRate()

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top: unit.Dp(30), Bottom: unit.Dp(30),
				Left: unit.Dp(30), Right: unit.Dp(30),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
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
						status := fmt.Sprintf("%d / %d", our_height, best_height)
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
						status := fmt.Sprintf("%d", peer_count)
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
						status := utils.FormatHashRate(network_hash_rate)
						label := material.Label(th, unit.Sp(22), status)
						label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
						return label.Layout(gtx)
					}),

					layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Label(th, unit.Sp(18), "TXp / Time Sync")
						label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 100}
						return label.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						status := fmt.Sprintf("%d:%d / %s | %s | %s",
							mempool_tx_count,
							regpool_tx_count,
							globals.GetOffset().Round(time.Millisecond).String(),
							globals.GetOffsetNTP().Round(time.Millisecond).String(),
							globals.GetOffsetP2P().Round(time.Millisecond).String(),
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
						value := utils.FormatBytes(p.nodeSize.size)
						label := material.Label(th, unit.Sp(22), value)
						label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
						return label.Layout(gtx)
					}),

					layout.Rigid(layout.Spacer{Height: unit.Dp(30)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.X = gtx.Dp(150)
						return p.buttonShowCommand.Layout(gtx, th)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return app_instance.Current.BottomBar.Layout(gtx, th)
		}),
	)
}

type NodeSize struct {
	size   int64
	active bool
}

func NewNodeSize() *NodeSize {
	nodeSize := &NodeSize{
		active: false,
		size:   0,
	}

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			if nodeSize.active {
				err := nodeSize.Calculate()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	err := nodeSize.Calculate()
	if err != nil {
		log.Fatal(err)
	}
	return nodeSize
}

func (n *NodeSize) Calculate() error {
	nodeDir := app_instance.Current.Settings.NodeDir
	size, err := utils.GetFolderSize(nodeDir)
	if err != nil {
		return err
	}

	n.size = size
	return nil
}
