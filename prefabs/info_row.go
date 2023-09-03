package prefabs

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/theme"
)

type InfoRow struct {
	Editor *widget.Editor
}

func NewInfoRow() *InfoRow {
	return &InfoRow{
		Editor: &widget.Editor{ReadOnly: true},
	}
}

func (i *InfoRow) Layout(gtx layout.Context, th *material.Theme, title string, value string) layout.Dimensions {
	if i.Editor.Text() != value {
		i.Editor.SetText(value)
	}

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), title)
			lbl.Font.Weight = font.Bold
			lbl.Color = theme.Current.TextMuteColor
			return lbl.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			editor := material.Editor(th, i.Editor, "")
			return editor.Layout(gtx)
		}),
	)
}
