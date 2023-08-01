package page_wallet

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/rpc"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/utils"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type PageTransaction struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation
	entry          *rpc.Entry
	txIdEditor     *widget.Editor

	srcImgCoinbase paint.ImageOp
	srcImgDown     paint.ImageOp
	srcImgUp       paint.ImageOp
	txTypeImg      components.Image

	list *widget.List
}

var _ router.Page = &PageTransaction{}

func NewPageTransaction() *PageTransaction {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	txIdEditor := new(widget.Editor)
	txIdEditor.WrapPolicy = text.WrapGraphemes
	txIdEditor.ReadOnly = true

	imgUp, _ := assets.GetImage("arrow_up_arc.png")
	srcImgUp := paint.NewImageOp(imgUp)

	imgDown, _ := assets.GetImage("arrow_down_arc.png")
	srcImgDown := paint.NewImageOp(imgDown)

	imgCoinbase, _ := assets.GetImage("coinbase.png")
	srcImgCoinbase := paint.NewImageOp(imgCoinbase)

	txTypeImg := components.Image{
		Fit: components.Cover,
	}

	list := new(widget.List)
	list.Axis = layout.Vertical

	return &PageTransaction{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		list:       list,
		txIdEditor: txIdEditor,

		srcImgCoinbase: srcImgCoinbase,
		srcImgDown:     srcImgDown,
		srcImgUp:       srcImgUp,
		txTypeImg:      txTypeImg,
	}
}

func (p *PageTransaction) IsActive() bool {
	return p.isActive
}

func (p *PageTransaction) Enter() {
	p.txIdEditor.SetText(p.entry.TXID)
	fmt.Println(p.entry)
	if p.entry.Incoming {
		p.txTypeImg.Src = p.srcImgDown
	} else {
		p.txTypeImg.Src = p.srcImgUp
	}

	page_instance.header.SetTitle(lang.Translate("Transaction"))
	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		txId := utils.ReduceTxId(p.entry.TXID)
		if txId == "" {
			txId = "From Coinbase"
		}

		lbl := material.Label(th, unit.Sp(16), txId)
		return lbl.Layout(gtx)
	}

	p.isActive = true
	if !page_instance.header.IsHistory(PAGE_TRANSACTION) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageTransaction) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageTransaction) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.entry.TXID != "" {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					r := op.Record(gtx.Ops)
					dims := layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis:      layout.Horizontal,
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Max.X = gtx.Dp(50)
								gtx.Constraints.Max.Y = gtx.Dp(50)
								return p.txTypeImg.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								editor := material.Editor(th, p.txIdEditor, "")
								return editor.Layout(gtx)
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
			)
		})
	}

	if !p.entry.Coinbase {
		if p.entry.Incoming {
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), lang.Translate("Sender"))
						lbl.Font.Weight = font.Bold
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), p.entry.Sender)
						return lbl.Layout(gtx)
					}),
				)
			})
		} else {
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), lang.Translate("Destination"))
						lbl.Font.Weight = font.Bold
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), p.entry.Destination)
						return lbl.Layout(gtx)
					}),
				)
			})
		}
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		amount := globals.FormatMoney(p.entry.Amount)
		fees := globals.FormatMoney(p.entry.Fees)
		burn := globals.FormatMoney(p.entry.Burn)
		date := p.entry.Time.Format("2006-01-02 15:04")
		timeAgo := lang.TimeAgo(p.entry.Time)

		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			InfoRowLayout{}.Layout(gtx, th, lang.Translate("Amount"), amount),
			InfoRowLayout{}.Layout(gtx, th, lang.Translate("Fees"), fees),
			InfoRowLayout{}.Layout(gtx, th, lang.Translate("Burn"), burn),
			InfoRowLayout{}.Layout(gtx, th, lang.Translate("Block Height"), fmt.Sprint(p.entry.Height)),
			InfoRowLayout{}.Layout(gtx, th, lang.Translate("Date"), date),
			InfoRowLayout{}.Layout(gtx, th, lang.Translate("Time"), timeAgo),
		)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(th, unit.Sp(16), lang.Translate("Block"))
				lbl.Font.Weight = font.Bold
				return lbl.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(th, unit.Sp(16), p.entry.BlockHash)
				return lbl.Layout(gtx)
			}),
		)
	})

	if !p.entry.Coinbase {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), lang.Translate("Proof"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), p.entry.Proof)
					return lbl.Layout(gtx)
				}),
			)
		})
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
	})

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

type InfoRowLayout struct {
	editor *widget.Editor
}

func (i InfoRowLayout) Layout(gtx layout.Context, th *material.Theme, title string, value string) layout.FlexChild {
	if i.editor == nil {
		i.editor = new(widget.Editor)
		i.editor.SetText(value)
	}

	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(th, unit.Sp(16), title)
				lbl.Font.Weight = font.Bold
				return lbl.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				editor := material.Editor(th, i.editor, "")
				return editor.Layout(gtx)
			}),
		)
	})
}
