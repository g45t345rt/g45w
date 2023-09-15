package page_wallet

import (
	"fmt"
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/rpc"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type PageTransaction struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation
	entry          *rpc.Entry

	txTypeImg components.Image

	txIdEditor              *widget.Editor
	senderDestinationEditor *widget.Editor
	blockHashEditor         *widget.Editor
	proofEditor             *widget.Editor
	infoRows                []*prefabs.InfoRow

	payloadArgInfoList []PayloadArgInfo

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

	txTypeImg := components.Image{
		Fit: components.Cover,
	}

	list := new(widget.List)
	list.Axis = layout.Vertical

	return &PageTransaction{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		list:                    list,
		txIdEditor:              &widget.Editor{ReadOnly: true},
		senderDestinationEditor: &widget.Editor{ReadOnly: true},
		blockHashEditor:         &widget.Editor{ReadOnly: true},
		proofEditor:             &widget.Editor{ReadOnly: true},
		infoRows:                prefabs.NewInfoRows(6),

		txTypeImg: txTypeImg,
	}
}

func (p *PageTransaction) IsActive() bool {
	return p.isActive
}

func (p *PageTransaction) Clear() {
	p.payloadArgInfoList = make([]PayloadArgInfo, 0)
	p.txIdEditor.SetText("")
	p.senderDestinationEditor.SetText("")
	p.blockHashEditor.SetText("")
	p.proofEditor.SetText("")
}

func (p *PageTransaction) Enter() {
	p.Clear()
	p.txIdEditor.SetText(p.entry.TXID)

	for _, arg := range p.entry.Payload_RPC {
		p.payloadArgInfoList = append(p.payloadArgInfoList, *NewPayloadArgLayout(arg))
	}

	if p.entry.Incoming {
		sender := p.entry.Sender
		if sender == "" {
			sender = "?"
		}

		p.senderDestinationEditor.SetText(sender)
		p.txTypeImg.Src = theme.Current.ArrowDownArcImage
	} else {
		p.senderDestinationEditor.SetText(p.entry.Destination)
		p.txTypeImg.Src = theme.Current.ArrowUpArcImage
	}

	p.blockHashEditor.SetText(p.entry.BlockHash)
	p.proofEditor.SetText(p.entry.Proof)

	page_instance.header.Title = func() string { return lang.Translate("Transaction") }
	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		txId := utils.ReduceTxId(p.entry.TXID)
		if txId == "" {
			txId = lang.Translate("From Coinbase")
		}

		lbl := material.Label(th, unit.Sp(16), txId)
		lbl.Color = theme.Current.TextMuteColor
		return lbl.Layout(gtx)
	}
	page_instance.header.ButtonRight = nil

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

					paint.FillShape(gtx.Ops, theme.Current.ListBgColor,
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
						lbl.Color = theme.Current.TextMuteColor
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						editor := material.Editor(th, p.senderDestinationEditor, "")
						return editor.Layout(gtx)
					}),
				)
			})
		} else {
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), lang.Translate("Destination"))
						lbl.Font.Weight = font.Bold
						lbl.Color = theme.Current.TextMuteColor
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						editor := material.Editor(th, p.senderDestinationEditor, "")
						return editor.Layout(gtx)
					}),
				)
			})
		}
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				amount := globals.FormatMoney(p.entry.Amount)
				return p.infoRows[0].Layout(gtx, th, lang.Translate("Amount"), amount)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				fees := globals.FormatMoney(p.entry.Fees)
				return p.infoRows[1].Layout(gtx, th, lang.Translate("Fees"), fees)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				burn := globals.FormatMoney(p.entry.Burn)
				return p.infoRows[2].Layout(gtx, th, lang.Translate("Burn"), burn)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				blockHeight := fmt.Sprint(p.entry.Height)
				return p.infoRows[3].Layout(gtx, th, lang.Translate("Block Height"), blockHeight)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				date := p.entry.Time.Format("2006-01-02 15:04")
				return p.infoRows[4].Layout(gtx, th, lang.Translate("Date"), date)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				timeAgo := lang.TimeAgo(p.entry.Time)
				return p.infoRows[5].Layout(gtx, th, lang.Translate("Time"), timeAgo)
			}),
		)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(th, unit.Sp(16), lang.Translate("Block Hash"))
				lbl.Font.Weight = font.Bold
				lbl.Color = theme.Current.TextMuteColor
				return lbl.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				editor := material.Editor(th, p.blockHashEditor, "")
				return editor.Layout(gtx)
			}),
		)
	})

	if !p.entry.Coinbase {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), lang.Translate("Proof"))
					lbl.Font.Weight = font.Bold
					lbl.Color = theme.Current.TextMuteColor
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(th, p.proofEditor, "")
					return editor.Layout(gtx)
				}),
			)
		})
	}

	for i := range p.payloadArgInfoList {
		idx := i
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return p.payloadArgInfoList[idx].Layout(gtx, th)
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

type PayloadArgInfo struct {
	arg    rpc.Argument
	editor *widget.Editor
}

func NewPayloadArgLayout(arg rpc.Argument) *PayloadArgInfo {
	editor := new(widget.Editor)
	editor.ReadOnly = true
	editor.SetText(fmt.Sprint(arg.Value))

	return &PayloadArgInfo{
		editor: editor,
		arg:    arg,
	}
}

func (p *PayloadArgInfo) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	var name string
	switch p.arg.Name {
	case "D":
		name = lang.Translate("Destination Port")
	case "S":
		name = lang.Translate("Source Port")
	case "V":
		name = lang.Translate("Value Transfer")
	case "C":
		name = lang.Translate("Comment")
	case "E":
		name = lang.Translate("Expiry")
	case "R":
		name = lang.Translate("Replyback Address")
	case "N":
		name = lang.Translate("Needs Replyback Address")
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), name)
			lbl.Font.Weight = font.Bold
			lbl.Color = theme.Current.TextMuteColor
			return lbl.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			editor := material.Editor(th, p.editor, "")
			return editor.Layout(gtx)
		}),
	)
}
