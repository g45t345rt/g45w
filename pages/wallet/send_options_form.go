package page_wallet

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type PageSendOptionsForm struct {
	isActive bool

	txtComment     *components.TextField
	txtDescription *components.TextField
	txtDstPort     *components.TextField

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	list *widget.List
}

var _ router.Page = &PageSendOptionsForm{}

func NewPageSendOptionsForm() *PageSendOptionsForm {
	txtComment := components.NewTextField()
	txtComment.Editor().SingleLine = false
	txtComment.Editor().Submit = false
	txtDescription := components.NewTextField()
	txtDescription.Editor().SingleLine = false
	txtDescription.Editor().Submit = false
	txtDstPort := components.NewNumberTextField()

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	return &PageSendOptionsForm{
		txtComment:     txtComment,
		txtDstPort:     txtDstPort,
		txtDescription: txtDescription,
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		list:           list,
	}
}

func (p *PageSendOptionsForm) IsActive() bool {
	return p.isActive
}

func (p *PageSendOptionsForm) Enter() {
	p.isActive = true
	if !page_instance.header.IsHistory(PAGE_SEND_OPTIONS_FORM) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
	page_instance.header.Title = func() string { return lang.Translate("Send Options") }
	page_instance.header.Subtitle = nil
}

func (p *PageSendOptionsForm) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageSendOptionsForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), lang.Translate("When using an integrated address, the options for \"Comment\" and \"Destination Port\" are discarded."))
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					p.txtComment.Input.EditorMinY = gtx.Dp(75)
					return p.txtComment.Layout(gtx, th, lang.Translate("Comment"), lang.Translate("The comment is store on the blockchain and natively encrypted. Only the sender / receiver can decrypt."))
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtDstPort.Layout(gtx, th, lang.Translate("Destination Port"), "Specific service port.")
		},
		func(gtx layout.Context) layout.Dimensions {
			p.txtDescription.Input.EditorMinY = gtx.Dp(75)
			return p.txtDescription.Layout(gtx, th, lang.Translate("Description"), lang.Translate("Saved locally in your wallet."))
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if p.txtComment.Input.Clickable.Clicked() {
		p.list.ScrollTo(0)
	}

	if p.txtDstPort.Input.Clickable.Clicked() {
		p.list.ScrollTo(1)
	}

	if p.txtDescription.Input.Clickable.Clicked() {
		p.list.ScrollTo(2)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}
