package page_wallet

import (
	"bytes"
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/wallet_manager"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type PageReceiveForm struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	list       *widget.List
	addrEditor *widget.Editor
	addrImage  *components.Image
}

var _ router.Page = &PageReceiveForm{}

func NewPageReceiveForm() *PageReceiveForm {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	addrEditor := new(widget.Editor)
	addrEditor.WrapPolicy = text.WrapGraphemes
	addrEditor.Alignment = text.Middle
	addrEditor.ReadOnly = true

	return &PageReceiveForm{
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		list:           list,
		addrEditor:     addrEditor,
	}
}

func (p *PageReceiveForm) IsActive() bool {
	return p.isActive
}

func (p *PageReceiveForm) Enter() {
	p.isActive = true

	if !page_instance.header.IsHistory(PAGE_RECEIVE_FORM) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
	page_instance.pageBalanceTokens.ResetWalletHeader()

	addr := wallet_manager.OpenedWallet.Info.Addr
	imgBytes, _ := qrcode.Encode(addr, qrcode.Medium, 256)
	img, _, _ := image.Decode(bytes.NewBuffer(imgBytes))

	p.addrImage = &components.Image{
		Src: paint.NewImageOp(img),
		Fit: components.Contain,
	}

	p.addrEditor.SetText(addr)
}

func (p *PageReceiveForm) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageReceiveForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = gtx.Dp(250)
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				editor := material.Editor(th, p.addrEditor, "")
				editor.TextSize = unit.Sp(16)
				editor.Font.Weight = font.Bold
				return editor.Layout(gtx)
			})
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.Y = gtx.Dp(250)
				return p.addrImage.Layout(gtx)
			})
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}
