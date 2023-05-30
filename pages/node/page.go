package page_node

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/components"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Page struct {
	isActive bool

	buttonShowCommand *components.Button
	hubIcon           *widget.Icon
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
	}
}

func (p *Page) IsActive() bool {
	return p.isActive
}

func (p *Page) Enter() {
	app_instance.Current.BottomBar.SetActive("node")
	p.isActive = true
}

func (p *Page) Leave() {
	p.isActive = false
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	paint.ColorOp{Color: color.NRGBA{A: 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

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
						label := material.Label(th, unit.Sp(22), "43534 / 34565")
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
						label := material.Label(th, unit.Sp(22), "105")
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
						label := material.Label(th, unit.Sp(22), "380 MH/s")
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
						label := material.Label(th, unit.Sp(22), "4:0 / 0s | 0s | 64ms")
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
						label := material.Label(th, unit.Sp(22), "523 MB")
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
