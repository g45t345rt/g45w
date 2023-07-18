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
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type NotificationStyle struct {
	BgColor    color.NRGBA
	TextColor  color.NRGBA
	Icon       *widget.Icon
	Direction  layout.Direction
	OuterInset layout.Inset
	InnerInset layout.Inset
	Rounded    Rounded
	Animation  ModalAnimation
}

type NotificationModal struct {
	Style NotificationStyle
	Modal *Modal

	window   *app.Window
	title    string
	subtitle string
	timer    *time.Timer
}

func NewNotificationErrorModal(window *app.Window) *NotificationModal {
	iconError, _ := widget.NewIcon(icons.AlertError)
	return NewNotificationModal(window,
		NotificationStyle{
			BgColor:    color.NRGBA{R: 255, A: 255},
			TextColor:  color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			Direction:  layout.N,
			OuterInset: layout.UniformInset(unit.Dp(10)),
			InnerInset: layout.Inset{
				Top: unit.Dp(10), Bottom: unit.Dp(10),
				Left: unit.Dp(15), Right: unit.Dp(15),
			},
			Rounded:   UniformRounded(unit.Dp(10)),
			Icon:      iconError,
			Animation: NewModalAnimationDown(),
		},
	)
}

func NewNotificationSuccessModal(window *app.Window) *NotificationModal {
	iconSuccess, _ := widget.NewIcon(icons.ActionCheckCircle)
	return NewNotificationModal(window,
		NotificationStyle{
			BgColor:    color.NRGBA{R: 0, G: 255, B: 0, A: 255},
			TextColor:  color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			Direction:  layout.N,
			OuterInset: layout.UniformInset(unit.Dp(10)),
			InnerInset: layout.Inset{
				Top: unit.Dp(10), Bottom: unit.Dp(10),
				Left: unit.Dp(15), Right: unit.Dp(15),
			},
			Rounded:   UniformRounded(unit.Dp(10)),
			Icon:      iconSuccess,
			Animation: NewModalAnimationDown(),
		},
	)
}

func NewNotificationInfoModal(window *app.Window) *NotificationModal {
	iconInfo, _ := widget.NewIcon(icons.ActionInfo)
	return NewNotificationModal(window,
		NotificationStyle{
			BgColor:    color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			TextColor:  color.NRGBA{A: 255},
			Direction:  layout.N,
			OuterInset: layout.UniformInset(unit.Dp(10)),
			InnerInset: layout.Inset{
				Top: unit.Dp(10), Bottom: unit.Dp(10),
				Left: unit.Dp(15), Right: unit.Dp(15),
			},
			Rounded:   UniformRounded(unit.Dp(10)),
			Icon:      iconInfo,
			Animation: NewModalAnimationDown(),
		},
	)
}

func NewNotificationModal(window *app.Window, style NotificationStyle) *NotificationModal {
	modalStyle := ModalStyle{
		CloseOnOutsideClick: false,
		CloseOnInsideClick:  true,
		BgColor:             style.BgColor,
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
	return n.Modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return n.Style.InnerInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
	})
}
