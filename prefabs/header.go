package prefabs

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/olebedev/emitter"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Header struct {
	LabelTitle  material.LabelStyle
	ButtonRight *components.Button

	buttonBack *components.Button
	router     *router.Router

	history []interface{}
}

func NewHeader(labelTitle material.LabelStyle, r *router.Router, buttonRight *components.Button) *Header {
	textColor := color.NRGBA{R: 0, G: 0, B: 0, A: 100}
	textHoverColor := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

	walletIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	buttonBack := components.NewButton(components.ButtonStyle{
		Icon:           walletIcon,
		TextColor:      textColor,
		HoverTextColor: &textHoverColor,
	})

	header := &Header{
		LabelTitle:  labelTitle,
		buttonBack:  buttonBack,
		router:      r,
		ButtonRight: buttonRight,
		history:     make([]interface{}, 0),
	}

	r.Event.On(router.EVENT_CHANGE, func(e *emitter.Event) {
		header.history = append(header.history, e.Args[0])
	})

	return header
}

func (h *Header) SetTitle(title string) {
	h.LabelTitle.Text = title
}

func (h *Header) back() {
	tag := h.history[len(h.history)-2]
	h.router.SetCurrent(tag)
	h.history = h.history[:len(h.history)-2]
}

func (h *Header) Layout(gtx layout.Context, th *material.Theme, subWidget layout.Widget) layout.Dimensions {
	if h.buttonBack.Clickable.Clicked() {
		h.back()
	}

	showBackButton := len(h.history) > 1

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if showBackButton {
				gtx.Constraints.Min.X = gtx.Dp(30)
				gtx.Constraints.Min.Y = gtx.Dp(30)
			} else {
				gtx.Constraints.Max.X = 0
				gtx.Constraints.Max.Y = 0
			}
			return h.buttonBack.Layout(gtx, th)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if showBackButton {
				return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
			}

			return layout.Dimensions{}
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return h.LabelTitle.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if subWidget != nil {
						return subWidget(gtx)
					}
					return layout.Dimensions{}
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !showBackButton && h.ButtonRight != nil {
				gtx.Constraints.Max.X = gtx.Dp(25)
				gtx.Constraints.Max.Y = gtx.Dp(25)
				return h.ButtonRight.Layout(gtx, th)
			}

			return layout.Dimensions{}
		}),
	)
}
