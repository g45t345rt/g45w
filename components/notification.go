package components

import (
	"image/color"
	"time"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type NotificationColors struct {
	TextColor       color.NRGBA
	BackgroundColor color.NRGBA
}

type NotificationStyle struct {
	Icon       *widget.Icon
	Direction  layout.Direction
	OuterInset layout.Inset
	InnerInset layout.Inset
	Rounded    Rounded
	Animation  ModalAnimation
	Colors     NotificationColors
}

type NotificationModal struct {
	Style NotificationStyle
	Modal *Modal

	window   *app.Window
	title    string
	subtitle string
	timer    *time.Timer
}

func NewNotificationModal(window *app.Window, style NotificationStyle) *NotificationModal {
	modalStyle := ModalStyle{
		CloseOnOutsideClick: false,
		CloseOnInsideClick:  true,
		Rounded:             style.Rounded,
		Direction:           style.Direction,
		Inset:               style.OuterInset,
		Animation:           style.Animation,
	}

	modal := NewModal(modalStyle)
	notification := &NotificationModal{
		window: window,
		Style:  style,
		Modal:  modal,
	}
	return notification
}

func (n *NotificationModal) SetText(title string, subtitle string) {
	n.title = title
	n.subtitle = subtitle
}

func (n *NotificationModal) SetVisible(visible bool, closeAfter time.Duration) {
	if visible {
		if n.timer != nil {
			n.timer.Stop()
		}

		if closeAfter > 0 {
			n.timer = time.AfterFunc(closeAfter, func() {
				n.SetVisible(false, 0)
				n.window.Invalidate()
			})
		}
	}

	n.Modal.SetVisible(visible)
}

func (n *NotificationModal) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	textColor := n.Style.Colors.TextColor
	n.Modal.Style.Colors.BackgroundColor = n.Style.Colors.BackgroundColor
	return n.Modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return n.Style.InnerInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if n.Style.Icon != nil {
						return n.Style.Icon.Layout(gtx, textColor)
					}

					return layout.Dimensions{}
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Label(th, unit.Sp(18), n.title)
							label.Font.Weight = font.Bold
							label.Color = textColor
							return label.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Label(th, unit.Sp(16), n.subtitle)
							label.Color = textColor
							return label.Layout(gtx)
						}),
					)
				}),
			)
		})
	})
}
