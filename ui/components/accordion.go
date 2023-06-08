package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type AccordionStyle struct {
	Button *Button
	Border widget.Border
	Inset  layout.Inset
}

type Accordion struct {
	Style   AccordionStyle
	Visible bool
}

func NewAccordion(style AccordionStyle, visible bool) *Accordion {
	return &Accordion{
		Visible: visible,
		Style:   style,
	}
}

func (a *Accordion) Layout(gtx layout.Context, th *material.Theme, w layout.Widget) layout.Dimensions {
	if a.Style.Button.Clickable.Clicked() {
		a.Visible = !a.Visible
	}

	return a.Style.Border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return a.Style.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.Style.Button.Layout(gtx, th)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if a.Visible {
						return w(gtx)
					}

					return layout.Dimensions{}
				}),
			)
		})
	})
}
