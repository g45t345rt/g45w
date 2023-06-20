package page_wallet

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/containers/bottom_bar"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/utils"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageTxs struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	txItems []*TxListItem

	listStyle material.ListStyle
}

var _ router.Page = &PageTxs{}

func NewPageTxs() *PageTxs {
	th := app_instance.Theme

	txItems := []*TxListItem{}

	upIcon, _ := widget.NewIcon(icons.NavigationArrowUpward)
	downIcon, _ := widget.NewIcon(icons.NavigationArrowDownward)

	for i := 0; i < 10; i++ {
		icon := upIcon
		if i%2 == 0 {
			icon = downIcon
		}

		txItems = append(txItems, &TxListItem{
			asset:     fmt.Sprintf("DERO (%s)", utils.ReduceString(crypto.ZEROHASH.String(), 4, 4)),
			date:      time.Now().Format("2006-01-02"),
			addr:      utils.ReduceAddr("dero1qy4egwhfjxtt96pdlkkc3d99yqkewkvljf50unqh5vz7xepl8w6yyqgklpjsz"),
			amount:    "342.35546",
			Clickable: new(widget.Clickable),
			Icon:      icon,
		})
	}

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical
	listStyle := material.List(th, list)
	listStyle.AnchorStrategy = material.Overlay

	return &PageTxs{
		txItems:        txItems,
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		listStyle:      listStyle,
	}
}

func (p *PageTxs) IsActive() bool {
	return p.isActive
}

func (p *PageTxs) Enter() {
	p.isActive = true

	bottom_bar.Instance.SetButtonActive(bottom_bar.BUTTON_TXS)
	p.animationEnter.Start()
	p.animationLeave.Reset()
}

func (p *PageTxs) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageTxs) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	{
		state := p.animationEnter.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	{
		state := p.animationLeave.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}

		if state.Finished {
			p.isActive = false
			op.InvalidateOp{}.Add(gtx.Ops)
		}
	}

	widgets := []layout.Widget{}

	for _, item := range p.txItems {
		widgets = append(widgets, item.Layout)
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
	})

	return p.listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return widgets[index](gtx)
	})
}

type TxListItem struct {
	asset  string
	addr   string
	amount string
	date   string

	Clickable *widget.Clickable
	Icon      *widget.Icon
}

func (item *TxListItem) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Inset{
		Top: unit.Dp(0), Bottom: unit.Dp(10),
		Right: unit.Dp(30), Left: unit.Dp(30),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		th := app_instance.Theme
		m := op.Record(gtx.Ops)
		dims := item.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top: unit.Dp(10), Bottom: unit.Dp(10),
				Left: unit.Dp(15), Right: unit.Dp(15),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.X = gtx.Dp(50)
						gtx.Constraints.Max.Y = gtx.Dp(50)
						return item.Icon.Layout(gtx, color.NRGBA{A: 255})
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis:      layout.Horizontal,
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										label := material.Label(th, unit.Sp(18), item.addr)
										label.Font.Weight = font.Bold
										return label.Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										label := material.Label(th, unit.Sp(14), item.date)
										label.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 150}
										return label.Layout(gtx)
									}),
								)
							}),
							layout.Flexed(1, layout.Spacer{}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										label := material.Label(th, unit.Sp(18), item.amount)
										label.Font.Weight = font.Bold
										return label.Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										label := material.Label(th, unit.Sp(14), item.asset)
										label.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 150}
										return label.Layout(gtx)
									}),
								)
							}),
						)
					}),
				)
			})
		})
		c := m.Stop()

		paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			clip.RRect{
				Rect: image.Rectangle{Max: dims.Size},
				SE:   gtx.Dp(10),
				NW:   gtx.Dp(10),
				NE:   gtx.Dp(10),
				SW:   gtx.Dp(10),
			}.Op(gtx.Ops))

		c.Add(gtx.Ops)

		return dims
	})
	/*
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return dims
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
		)


			if item.Clickable.Hovered() {
				pointer.CursorPointer.Add(gtx.Ops)
				bounds := image.Rect(0, 0, dims.Size.X, dims.Size.Y)
				paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 100},
					clip.UniformRRect(bounds, 10).Op(gtx.Ops),
				)
			}*/
}
