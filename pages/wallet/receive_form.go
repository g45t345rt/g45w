package page_wallet

import (
	"bytes"
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/wallet_manager"
	qrcode "github.com/skip2/go-qrcode"
)

type PageReceiveForm struct {
	isActive bool

	headerPageAnimation *prefabs.PageHeaderAnimation

	list       *widget.List
	addrEditor *widget.Editor
	addrImage  *components.Image
}

var _ router.Page = &PageReceiveForm{}

func NewPageReceiveForm() *PageReceiveForm {

	list := new(widget.List)
	list.Axis = layout.Vertical

	addrEditor := new(widget.Editor)
	addrEditor.WrapPolicy = text.WrapGraphemes
	addrEditor.Alignment = text.Middle
	addrEditor.ReadOnly = true

	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_RECEIVE_FORM)

	return &PageReceiveForm{
		headerPageAnimation: headerPageAnimation,
		list:                list,
		addrEditor:          addrEditor,
	}
}

func (p *PageReceiveForm) IsActive() bool {
	return p.isActive
}

func (p *PageReceiveForm) Enter() {
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)

	page_instance.pageBalanceTokens.ResetWalletHeader()

	addr := wallet_manager.OpenedWallet.Memory.GetAddress().String()
	imgBytes, _ := qrcode.Encode(addr, qrcode.Medium, 256)
	img, _, _ := image.Decode(bytes.NewBuffer(imgBytes))

	p.addrImage = &components.Image{
		Src: paint.NewImageOp(img),
		Fit: components.Contain,
		Rounded: components.Rounded{
			NW: unit.Dp(10), NE: unit.Dp(10),
			SW: unit.Dp(10), SE: unit.Dp(10),
		},
	}

	p.addrEditor.SetText(addr)
}

func (p *PageReceiveForm) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)
}

func (p *PageReceiveForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = gtx.Dp(260)
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				editor := material.Editor(th, p.addrEditor, "")
				editor.TextSize = unit.Sp(16)
				editor.Font.Weight = font.Bold
				return editor.Layout(gtx)
			})
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.Y = gtx.Dp(260)
				return p.addrImage.Layout(gtx, nil)
			})
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, widgets[index])
	})
}
