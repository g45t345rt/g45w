package components

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Input struct {
	TextSize    unit.Sp
	Editor      *widget.Editor
	EditorStyle material.EditorStyle
	EditorMinY  int
	Border      widget.Border
	Inset       layout.Inset

	Clickable widget.Clickable
}

func NewInput() *Input {
	editor := new(widget.Editor)
	editor.SingleLine = true
	editor.Submit = true
	editor.InputHint = key.HintText // Cap sentence flag
	border := widget.Border{Color: color.NRGBA{A: 240}, CornerRadius: unit.Dp(5), Width: unit.Dp(1)}

	return &Input{
		Editor: editor,
		Border: border,
		Inset: layout.Inset{
			Top: unit.Dp(15), Bottom: unit.Dp(15),
			Left: unit.Dp(12), Right: unit.Dp(12),
		},
	}
}

func NewPasswordInput() *Input {
	editor := new(widget.Editor)
	editor.SingleLine = true
	editor.Submit = true
	editor.InputHint = key.HintPassword
	editor.Mask = rune(42) // mask with *
	border := widget.Border{Color: color.NRGBA{A: 240}, CornerRadius: unit.Dp(5), Width: unit.Dp(1)}

	return &Input{
		Editor: editor,
		Border: border,
		Inset: layout.Inset{
			Top: unit.Dp(15), Bottom: unit.Dp(15),
			Left: unit.Dp(12), Right: unit.Dp(12),
		},
	}
}

func (t *Input) Value() string {
	return t.Editor.Text()
}

func (t *Input) SetValue(text string) {
	t.Editor.SetText(text)
}

func (t *Input) Layout(gtx layout.Context, th *material.Theme, hint string) layout.Dimensions {
	for _, e := range t.Editor.Events() {
		e, ok := e.(widget.SubmitEvent)
		if ok {
			fmt.Println(e.Text)
		}
	}

	gtx.Constraints.Min.Y = t.EditorMinY

	return t.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return t.Border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := t.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				t.EditorStyle = material.Editor(th, t.Editor, hint)
				if t.TextSize != 0 {
					t.EditorStyle.TextSize = t.TextSize
				}
				return t.EditorStyle.Layout(gtx)
			})
			call := macro.Stop()

			paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.UniformRRect(
				image.Rectangle{Max: dims.Size},
				int(t.Border.CornerRadius),
			).Op(gtx.Ops))

			call.Add(gtx.Ops)
			return dims
		})
	})
}
