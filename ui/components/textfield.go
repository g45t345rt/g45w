package components

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type TextField struct {
	TitleStyle material.LabelStyle
	Input      *Input
}

func NewTextField(th *material.Theme, title string, hint string) *TextField {
	titleStyle := material.Label(th, unit.Sp(20), title)
	titleStyle.Font.Weight = font.Bold
	input := NewInput(th, hint)

	return &TextField{
		TitleStyle: titleStyle,
		Input:      input,
	}
}

func NewPasswordTextField(th *material.Theme, title string, hint string) *TextField {
	titleStyle := material.Label(th, unit.Sp(20), title)
	titleStyle.Font.Weight = font.Bold

	input := NewPasswordInput(th, hint)

	return &TextField{
		TitleStyle: titleStyle,
		Input:      input,
	}
}

func (t *TextField) Value() string {
	return t.Input.Value()
}

func (t *TextField) SetValue(text string) {
	t.Input.SetValue(text)
}

func (t *TextField) Editor() *widget.Editor {
	return t.Input.Editor()
}

func (t *TextField) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.TitleStyle.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.Input.Layout(gtx, th)
		}),
	)
}
