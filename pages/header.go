package pages

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/components"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Header struct {
	LabelTitle  material.LabelStyle
	buttonBack  *components.Button
	childRouter *router.Router
}

func NewHeader(labelTitle material.LabelStyle, childRouter *router.Router) *Header {
	textColor := color.NRGBA{R: 0, G: 0, B: 0, A: 100}
	textHoverColor := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

	walletIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	buttonBack := components.NewButton(components.ButtonStyle{
		Icon:           walletIcon,
		TextColor:      textColor,
		HoverTextColor: &textHoverColor,
	})

	return &Header{
		LabelTitle:  labelTitle,
		buttonBack:  buttonBack,
		childRouter: childRouter,
	}
}

func (h *Header) Layout(gtx layout.Context, th *material.Theme, subWidget layout.Widget) layout.Dimensions {
	if h.buttonBack.Clickable.Clicked() {
		h.childRouter.SetCurrent(h.childRouter.Primary)
	}

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !h.childRouter.IsPrimary() {
				gtx.Constraints.Min.X = gtx.Dp(30)
				gtx.Constraints.Min.Y = gtx.Dp(30)
			} else {
				gtx.Constraints.Max.X = 0
				gtx.Constraints.Max.Y = 0
			}
			return h.buttonBack.Layout(gtx, th)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !h.childRouter.IsPrimary() {
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
	)
}
