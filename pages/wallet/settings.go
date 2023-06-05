package page_wallet

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
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

type PageSettings struct {
	isActive bool

	buttonDeleteWallet *components.Button
	confirmDelete      *components.Confirm

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	listStyle material.ListStyle
}

var _ router.Container = &PageSettings{}

func NewPageSettings() *PageSettings {
	th := app_instance.Current.Theme
	deleteIcon, _ := widget.NewIcon(icons.ActionDelete)
	buttonDeleteWallet := components.NewButton(components.ButtonStyle{
		Rounded:         unit.Dp(5),
		Text:            "DELETE WALLET",
		Icon:            deleteIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 255, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonDeleteWallet.Label.Alignment = text.Middle
	buttonDeleteWallet.Style.Font.Weight = font.Bold

	w := app_instance.Current.Window
	confirmDelete := components.NewConfirm(w, "Are you sure?", th, layout.Center)

	router := app_instance.Current.Router
	router.PushLayout(func(gtx layout.Context, th *material.Theme) {
		confirmDelete.Layout(gtx, th)
	})

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

	return &PageSettings{
		buttonDeleteWallet: buttonDeleteWallet,
		animationEnter:     animationEnter,
		animationLeave:     animationLeave,
		listStyle:          listStyle,
		confirmDelete:      confirmDelete,
	}
}

func (p *PageSettings) IsActive() bool {
	return p.isActive
}

func (p *PageSettings) Enter() {
	p.isActive = true
	p.animationEnter.Start()
	p.animationLeave.Reset()
}

func (p *PageSettings) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageSettings) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if p.buttonDeleteWallet.Clickable.Clicked() {
		p.confirmDelete.SetVisible(true)
	}

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
			return p.buttonDeleteWallet.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
		},
	}

	return p.listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}
