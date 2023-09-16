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

	window     *app.Window
	title      string
	text       string
	textEditor *widget.Editor
	timer      *time.Timer
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

	textEditor := new(widget.Editor)
	textEditor.ReadOnly = true

	modal := NewModal(modalStyle)
	notification := &NotificationModal{
		window:     window,
		Style:      style,
		Modal:      modal,
		textEditor: textEditor,
	}
	return notification
}

func (n *NotificationModal) SetText(title string, text string) {
	n.title = title
	n.text = text
}

func (n *NotificationModal) SetVisible(visible bool, closeAfter time.Duration) {
	if visible {
		if n.timer != nil {
			n.timer.Stop()
		}

		if closeAfter > 0 {
			n.timer = time.AfterFunc(closeAfter, func() {
				n.Modal.SetVisible(false)
				n.window.Invalidate()
			})
		}
	}

	n.Modal.SetVisible(visible)
	n.window.Invalidate()
}

func (n *NotificationModal) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	textColor := n.Style.Colors.TextColor
	n.Modal.Style.Colors.BackgroundColor = n.Style.Colors.BackgroundColor
	return n.Modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return n.Style.InnerInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
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
							editor := material.Editor(th, n.textEditor, "")
							editor.Color = textColor
							if n.textEditor.Text() != n.text {
								// using SetText here to avoid nil pointer while using SetText in the other func ???
								n.textEditor.SetText(n.text)
							}

							gtx.Constraints.Max.Y = gtx.Dp(150)
							return editor.Layout(gtx)
						}),
					)
				}),
			)
		})
	})
}
