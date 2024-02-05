package page_wallet

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageSendOptionsForm struct {
	isActive bool

	txtComment     *prefabs.TextField
	txtDescription *prefabs.TextField
	txtDstPort     *prefabs.TextField
	buttonContinue *components.Button

	headerPageAnimation *prefabs.PageHeaderAnimation

	list *widget.List
}

var _ router.Page = &PageSendOptionsForm{}

func NewPageSendOptionsForm() *PageSendOptionsForm {
	txtComment := prefabs.NewTextField()
	txtComment.Editor().SingleLine = false
	txtComment.Editor().Submit = false
	txtDescription := prefabs.NewTextField()
	txtDescription.Editor().SingleLine = false
	txtDescription.Editor().Submit = false
	txtDstPort := prefabs.NewNumberTextField()

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

	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_SEND_OPTIONS_FORM)
	return &PageSendOptionsForm{
		txtComment:          txtComment,
		txtDstPort:          txtDstPort,
		txtDescription:      txtDescription,
		headerPageAnimation: headerPageAnimation,
		list:                list,
		buttonContinue:      buttonContinue,
	}
}

func (p *PageSendOptionsForm) IsActive() bool {
	return p.isActive
}

func (p *PageSendOptionsForm) Enter() {
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)

	page_instance.header.Title = func() string { return lang.Translate("Send Options") }
	page_instance.header.Subtitle = nil
	page_instance.header.LeftLayout = nil
}

func (p *PageSendOptionsForm) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)
}

func (p *PageSendOptionsForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	if p.buttonContinue.Clicked(gtx) {
		page_instance.header.GoBack()
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), lang.Translate("When using an integrated address, the options for \"Comment\" and \"Destination Port\" are discarded."))
					lbl.Color = theme.Current.TextMuteColor
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					p.txtComment.Input.EditorMinY = gtx.Dp(75)
					return p.txtComment.Layout(gtx, th, lang.Translate("Comment"), lang.Translate("The comment is stored on the blockchain and natively encrypted. Only the sender / receiver can decrypt."))
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtDstPort.Layout(gtx, th, lang.Translate("Destination Port"), lang.Translate("Specific service port."))
		},
		func(gtx layout.Context) layout.Dimensions {
			p.txtDescription.Input.EditorMinY = gtx.Dp(75)
			return p.txtDescription.Layout(gtx, th, lang.Translate("Description"), lang.Translate("Saved locally in your wallet."))
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

	if p.txtComment.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(0)
	}

	if p.txtDstPort.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(1)
	}

	if p.txtDescription.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(2)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, widgets[index])
	})
}
