package components

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type TextField struct {
	TitleStyle  material.LabelStyle
	Hint        string
	EditorStyle material.EditorStyle
	EditorMinY  int
	Border      widget.Border

	clickable *widget.Clickable
}

func NewTextField(th *material.Theme, title string, hint string) *TextField {
	titleStyle := material.Label(th, unit.Sp(20), title)
	titleStyle.Font.Weight = font.Bold

	editor := new(widget.Editor)
	editor.SingleLine = true
	editor.Submit = true
	editorStyle := material.Editor(th, editor, hint)
	editorStyle.TextSize = unit.Sp(18)
	border := widget.Border{Color: color.NRGBA{A: 240}, CornerRadius: unit.Dp(5), Width: unit.Dp(1)}

	return &TextField{
		TitleStyle:  titleStyle,
		Hint:        hint,
		EditorStyle: editorStyle,
		clickable:   new(widget.Clickable),
		Border:      border,
	}
}

func NewPasswordTextField(th *material.Theme, title string, hint string) *TextField {
	titleStyle := material.Label(th, unit.Sp(20), title)
	titleStyle.Font.Weight = font.Bold

	editor := new(widget.Editor)
	editor.SingleLine = true
	editor.Submit = true
	editor.Mask = rune(42) // mask with *
	editorStyle := material.Editor(th, editor, hint)
	editorStyle.TextSize = unit.Sp(18)
	border := widget.Border{Color: color.NRGBA{A: 240}, CornerRadius: unit.Dp(5), Width: unit.Dp(1)}

	return &TextField{
		TitleStyle:  titleStyle,
		Hint:        hint,
		EditorStyle: editorStyle,
		Border:      border,
	}
}

func (t *TextField) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for _, e := range t.EditorStyle.Editor.Events() {
		e, ok := e.(widget.SubmitEvent)
		if ok {
			fmt.Println(e.Text)
		}
	}

	if t.clickable.Clicked() {
		t.EditorStyle.Editor.Focus() // able to click within Inset padding
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.TitleStyle.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.Y = t.EditorMinY
			return t.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return t.Border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					m := op.Record(gtx.Ops)
					dims := layout.Inset{
						Top:    unit.Dp(15),
						Bottom: unit.Dp(15),
						Left:   unit.Dp(12),
						Right:  unit.Dp(12),
					}.Layout(gtx, t.EditorStyle.Layout)
					c := m.Stop()

					paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.UniformRRect(
						image.Rectangle{Max: dims.Size},
						int(t.Border.CornerRadius),
					).Op(gtx.Ops))

					c.Add(gtx.Ops)
					return dims
				})
			})
		}),
	)
}
