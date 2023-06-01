package components

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Confirm struct {
	Prompt string

	modal     *Modal
	buttonYes *Button
	buttonNo  *Button

	clickedYes bool
	clickedNo  bool
}

func NewConfirm(prompt string, th *material.Theme, direction layout.Direction) *Confirm {
	modal := NewModal(ModalStyle{
		CloseOnOutsideClick: true,
		CloseOnInsideClick:  false,
		Direction:           direction,
		Inset:               layout.UniformInset(unit.Dp(10)),
		Animation:           NewModalAnimationScaleBounce(),
		Background:          NewModalBackground(),
	})

	buttonYes := NewButton(ButtonStyle{
		Rounded:         unit.Dp(5),
		Text:            "YES",
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       NewButtonAnimationDefault(),
	})
	buttonYes.Label.Alignment = text.Middle
	buttonYes.Style.Font.Weight = font.Bold

	buttonNo := NewButton(ButtonStyle{
		Rounded:         unit.Dp(5),
		Text:            "NO",
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       NewButtonAnimationDefault(),
	})
	buttonNo.Label.Alignment = text.Middle
	buttonNo.Style.Font.Weight = font.Bold

	return &Confirm{
		Prompt:     prompt,
		modal:      modal,
		buttonYes:  buttonYes,
		buttonNo:   buttonNo,
		clickedYes: false,
		clickedNo:  false,
	}
}

func (c *Confirm) ClickedYes() bool {
	return c.clickedYes
}

func (c *Confirm) ClickedNo() bool {
	return c.clickedNo
}

func (c *Confirm) SetVisible(gtx layout.Context, visible bool) {
	c.modal.SetVisible(gtx, visible)
}

func (c *Confirm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	c.clickedYes = c.buttonYes.Clickable.Clicked()
	c.clickedNo = c.buttonNo.Clickable.Clicked()

	if c.clickedYes || c.clickedNo {
		c.SetVisible(gtx, false)
	}

	return c.modal.Layout(gtx, func(gtx layout.Context) {

	}, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.Label(th, unit.Sp(18), c.Prompt)
					return label.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return c.buttonYes.Layout(gtx, th)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return c.buttonNo.Layout(gtx, th)
						}),
					)
				}),
			)
		})
	})
}
