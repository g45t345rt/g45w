package page_wallet

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageSendOptionsForm struct {
	isActive bool

	txtComment     *prefabs.TextField
	txtDstPort     *prefabs.TextField
	buttonContinue *components.Button

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	list *widget.List
}

var _ router.Page = &PageSendOptionsForm{}

func NewPageSendOptionsForm() *PageSendOptionsForm {
	txtComment := prefabs.NewTextField()
	txtComment.Editor().SingleLine = false
	txtComment.Editor().Submit = false
	txtDstPort := prefabs.NewNumberTextField()

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	arrowIcon, _ := widget.NewIcon(icons.NavigationChevronRight)
	buttonContinue := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      arrowIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonContinue.Label.Alignment = text.Middle
	buttonContinue.Style.Font.Weight = font.Bold

	return &PageSendOptionsForm{
		txtComment:     txtComment,
		txtDstPort:     txtDstPort,
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		list:           list,
		buttonContinue: buttonContinue,
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

	if p.buttonContinue.Clicked() {
		page_instance.header.GoBack()
	}

	widgets := []layout.
		Widget{
		func(gtx layout.Context) layout.Dimensions {

			return layout.
				Flex{Axis: layout.Vertical}.
				Layout(gtx, layout.
					Rigid(
						func(gtx layout.Context) layout.
							Dimensions {

							lbl := material.
								Label(th, unit.Sp(14),
									lang.Translate("The message and dst port are encrypted."+
										"\n\nOnly the sender / receiver can decrypt."))

							lbl.Color = theme.
								Current.
								TextMuteColor

							return lbl.
								Layout(gtx)
						}),

					layout.Rigid(
						layout.Spacer{Height: unit.Dp(10)}.Layout),

					layout.Rigid(
						func(
							gtx layout.Context) layout.Dimensions {
							p.
								txtComment.
								Input.
								EditorMinY = gtx.Dp(75)

							return p.
								txtComment.
								Layout(
									gtx,
									th,
									lang.
										Translate("Message"),
									lang.
										Translate("ex. secret loves you"))
						}),
				)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtDstPort.Layout(gtx, th, lang.Translate("DST Port"), lang.Translate("ex. 1337"))
		},

		func(gtx layout.Context) layout.Dimensions {
			p.buttonContinue.Style.Colors = theme.Current.ButtonPrimaryColors
			p.buttonContinue.Text = lang.Translate("CONTINUE")
			return p.buttonContinue.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
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

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}
