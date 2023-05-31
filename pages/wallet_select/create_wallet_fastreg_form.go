package page_wallet_select

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageCreateWalletFastRegForm struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	listStyle material.ListStyle

	txtThreadCount *components.TextField
	buttonStart    *components.Button
}

var _ router.Container = &PageCreateWalletFastRegForm{}

func NewPageCreateWalletFastRegForm() *PageCreateWalletFastRegForm {
	th := app_instance.Current.Theme
	list := new(widget.List)
	list.Axis = layout.Vertical
	listStyle := material.List(th, list)
	listStyle.AnchorStrategy = material.Overlay

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .5, ease.OutCubic),
	))

	txtThreadCount := components.NewTextField(th, "Thread Count", "")

	buildIcon, _ := widget.NewIcon(icons.HardwareMemory)
	buttonStart := components.NewButton(components.ButtonStyle{
		Rounded:         unit.Dp(5),
		Text:            "START",
		Icon:            buildIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})

	return &PageCreateWalletFastRegForm{
		listStyle:      listStyle,
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		txtThreadCount: txtThreadCount,
		buttonStart:    buttonStart,
	}
}

func (p *PageCreateWalletFastRegForm) Enter() {
	page_instance.header.LabelTitle.Text = "Fast Registration"
	p.isActive = true
	p.animationEnter.Start()
	p.animationLeave.Reset()
}

func (p *PageCreateWalletFastRegForm) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageCreateWalletFastRegForm) IsActive() bool {
	return p.isActive
}

func (p *PageCreateWalletFastRegForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	{
		state := p.animationEnter.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	{
		state := p.animationLeave.Update(gtx)
		if state.Finished {
			p.isActive = false
			op.InvalidateOp{}.Add(gtx.Ops)
		}

		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	if p.buttonStart.Clickable.Clicked() {

	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.txtThreadCount.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			label := material.Label(th, unit.Sp(16), "Progress")
			label.Font.Weight = font.Bold
			return label.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.buttonStart.Layout(gtx, th)
		},
	}

	return p.listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}
