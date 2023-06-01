package components

import (
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
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type NotificationStyle struct {
	BgColor    color.NRGBA
	TextColor  color.NRGBA
	Icon       *widget.Icon
	Direction  layout.Direction
	OuterInset layout.Inset
	InnerInset layout.Inset
	Rounded    unit.Dp
	Animation  ModalAnimation
}

type NotificationModal struct {
	Style NotificationStyle

	title    string
	subtitle string
	modal    *Modal
}

func NewNotificationErrorModal() *NotificationModal {
	iconError, _ := widget.NewIcon(icons.AlertError)
	return NewNotificationModal(
		NotificationStyle{
			BgColor:    color.NRGBA{R: 255, A: 255},
			TextColor:  color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			Direction:  layout.N,
			OuterInset: layout.UniformInset(unit.Dp(10)),
			InnerInset: layout.Inset{
				Top: unit.Dp(10), Bottom: unit.Dp(10),
				Left: unit.Dp(15), Right: unit.Dp(15),
			},
			Rounded:   unit.Dp(10),
			Icon:      iconError,
			Animation: NewModalAnimationDown(),
		},
	)
}

func NewNotificationSuccessModal() *NotificationModal {
	iconSuccess, _ := widget.NewIcon(icons.ActionCheckCircle)
	return NewNotificationModal(
		NotificationStyle{
			BgColor:    color.NRGBA{R: 0, G: 255, B: 0, A: 255},
			TextColor:  color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			Direction:  layout.N,
			OuterInset: layout.UniformInset(unit.Dp(10)),
			InnerInset: layout.Inset{
				Top: unit.Dp(10), Bottom: unit.Dp(10),
				Left: unit.Dp(15), Right: unit.Dp(15),
			},
			Rounded:   unit.Dp(10),
			Icon:      iconSuccess,
			Animation: NewModalAnimationDown(),
		},
	)
}

func NewNotificationModal(style NotificationStyle) *NotificationModal {
	modalStyle := ModalStyle{
		CloseOnOutsideClick: false,
		CloseOnInsideClick:  true,
		Direction:           style.Direction,
		Inset:               style.OuterInset,
		Animation:           style.Animation,
	}

	modal := NewModal(modalStyle)

	return &NotificationModal{
		Style: style,
		modal: modal,
	}
}

func (n *NotificationModal) SetText(title string, subtitle string) {
	n.title = title
	n.subtitle = subtitle
}

func (n *NotificationModal) SetVisible(gtx layout.Context, visible bool) {
	n.modal.SetVisible(gtx, visible)
}

func (n *NotificationModal) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return n.modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		r := op.Record(gtx.Ops)
		dims := n.Style.InnerInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if n.Style.Icon != nil {
						return n.Style.Icon.Layout(gtx, n.Style.TextColor)
					}

					return layout.Dimensions{}
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Label(th, unit.Sp(18), n.title)
							label.Font.Weight = font.Bold
							label.Color = n.Style.TextColor
							return label.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Label(th, unit.Sp(16), n.subtitle)
							label.Color = n.Style.TextColor
							return label.Layout(gtx)
						}),
					)
				}),
			)
		})
		c := r.Stop()

		rounded := gtx.Dp(n.Style.Rounded)
		paint.FillShape(gtx.Ops, n.Style.BgColor, clip.RRect{
			Rect: image.Rectangle{Max: dims.Size},
			NW:   rounded, NE: rounded,
			SE: rounded, SW: rounded,
		}.Op(gtx.Ops))

		c.Add(gtx.Ops)
		return dims
	})
}
