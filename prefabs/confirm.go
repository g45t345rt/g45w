package prefabs

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/theme"
)

type ConfirmText struct {
	Prompt string
	Yes    string
	No     string
}

type Confirm struct {
	Modal *components.Modal

	buttonYes *components.Button
	buttonNo  *components.Button

	clickedYes bool
	clickedNo  bool
}

func NewConfirm(direction layout.Direction) *Confirm {
	modal := components.NewModal(components.ModalStyle{
		CloseOnOutsideClick: true,
		CloseOnInsideClick:  false,
		Direction:           direction,
		Rounded:             components.UniformRounded(unit.Dp(10)),
		Inset:               layout.UniformInset(unit.Dp(10)),
		Animation:           components.NewModalAnimationScaleBounce(),
	})

	buttonYes := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		TextSize:  unit.Sp(14),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonYes.Label.Alignment = text.Middle
	buttonYes.Style.Font.Weight = font.Bold

	buttonNo := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		TextSize:  unit.Sp(14),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonNo.Label.Alignment = text.Middle
	buttonNo.Style.Font.Weight = font.Bold

	return &Confirm{
		Modal:      modal,
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

func (c *Confirm) SetVisible(visible bool) {
	c.Modal.SetVisible(visible)
}

func (c *Confirm) Layout(gtx layout.Context, th *material.Theme, text ConfirmText) layout.Dimensions {
	c.clickedYes = c.buttonYes.Clicked()
	c.clickedNo = c.buttonNo.Clicked()

	if c.clickedYes || c.clickedNo {
		c.SetVisible(false)
	}

	var lblSize layout.Dimensions
	c.Modal.Style.Colors = theme.Current.ModalColors
	return c.Modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.Label(th, unit.Sp(18), text.Prompt)
					lblSize = label.Layout(gtx)
					return lblSize
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = lblSize.Size.X
					return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							c.buttonNo.Text = text.No
							c.buttonNo.Style.Colors = theme.Current.ButtonPrimaryColors
							return c.buttonNo.Layout(gtx, th)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							c.buttonYes.Text = text.Yes
							c.buttonYes.Style.Colors = theme.Current.ButtonPrimaryColors
							return c.buttonYes.Layout(gtx, th)
						}),
					)
				}),
			)
		})
	})
}
