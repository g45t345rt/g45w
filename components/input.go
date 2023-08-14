package components

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type InputColors struct {
	BorderColor     color.NRGBA
	BackgroundColor color.NRGBA
	TextColor       color.NRGBA
	HintColor       color.NRGBA
}

type Input struct {
	FontWeight font.Weight
	TextSize   unit.Sp
	Editor     *widget.Editor
	EditorMinY int
	Border     widget.Border
	Inset      layout.Inset
	Clickable  *widget.Clickable
	Colors     InputColors

	keyboardClick *widget.Clickable
	submitted     bool
	submitText    string
	activeSubmit  bool
}

func NewInput() *Input {
	editor := new(widget.Editor)
	editor.SingleLine = true
	editor.Submit = true
	editor.InputHint = key.HintText // Cap sentence flag
	border := widget.Border{
		CornerRadius: unit.Dp(5),
		Width:        unit.Dp(1),
	}

	return &Input{
		Editor:        editor,
		Border:        border,
		Clickable:     new(widget.Clickable),
		keyboardClick: new(widget.Clickable),
		Inset: layout.Inset{
			Top: unit.Dp(15), Bottom: unit.Dp(15),
			Left: unit.Dp(12), Right: unit.Dp(12),
		},
	}
}

func NewNumberInput() *Input {
	input := NewInput()
	input.Editor.Filter = "0123456789."
	input.Editor.InputHint = key.HintNumeric
	return input
}

func NewPasswordInput() *Input {
	editor := new(widget.Editor)
	editor.SingleLine = true
	editor.Submit = true
	editor.InputHint = key.HintPassword
	editor.Mask = rune(42) // mask with *
	border := widget.Border{
		Color:        color.NRGBA{A: 240},
		CornerRadius: unit.Dp(5),
		Width:        unit.Dp(1),
	}

	return &Input{
		Editor:        editor,
		Border:        border,
		Clickable:     new(widget.Clickable),
		keyboardClick: new(widget.Clickable),
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

func (t *Input) Submitted() (bool, string) {
	t.activeSubmit = true
	if t.submitted {
		t.submitted = false
		return true, t.submitText
	}

	return false, t.submitText
}

func (t *Input) Layout(gtx layout.Context, th *material.Theme, hint string) layout.Dimensions {
	if t.activeSubmit {
		for _, e := range t.Editor.Events() {
			e, ok := e.(widget.SubmitEvent)
			if ok {
				t.SetValue("")
				t.submitText = e.Text
				t.submitted = true
			}
		}
	}

	gtx.Constraints.Min.Y = t.EditorMinY

	if t.keyboardClick.Clicked() {
		// on mobile if the keyboard popups and the input lose focus it will automatically close the keyboard
		// so we have to manually force keyboard request to avoid this issue
		key.SoftKeyboardOp{Show: true}.Add(gtx.Ops)
	}

	return t.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return t.keyboardClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			t.Border.Color = t.Colors.BorderColor
			return t.Border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				macro := op.Record(gtx.Ops)
				dims := t.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					editorStyle := material.Editor(th, t.Editor, hint)
					editorStyle.Color = t.Colors.TextColor
					editorStyle.HintColor = t.Colors.HintColor
					editorStyle.TextSize = th.TextSize
					if t.TextSize != 0 {
						editorStyle.TextSize = t.TextSize
					}
					editorStyle.Font.Weight = t.FontWeight
					return editorStyle.Layout(gtx)
				})
				call := macro.Stop()

				paint.FillShape(gtx.Ops, t.Colors.BackgroundColor, clip.UniformRRect(
					image.Rectangle{Max: dims.Size},
					int(t.Border.CornerRadius),
				).Op(gtx.Ops))

				call.Add(gtx.Ops)
				return dims
			})

		})
	})
}
