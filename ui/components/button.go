package components

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type ButtonAnimation struct {
	animationEnter *animation.Animation
	transformEnter animation.TransformFunc
	animationLeave *animation.Animation
	transformLeave animation.TransformFunc
	animationClick *animation.Animation
	transformClick animation.TransformFunc
}

type ButtonStyle struct {
	TextColor            color.NRGBA
	BackgroundColor      color.NRGBA
	Rounded              unit.Dp
	Text                 string
	TextSize             unit.Sp
	Inset                layout.Inset
	Font                 font.Font
	Icon                 *widget.Icon
	IconGap              unit.Dp
	HoverBackgroundColor *color.NRGBA
	HoverTextColor       *color.NRGBA
	Animation            ButtonAnimation
}

type Button struct {
	Style     ButtonStyle
	Clickable *widget.Clickable
	Label     *widget.Label
	Focused   bool

	animClickable    *widget.Clickable
	hoverSwitchState bool
}

func NewButtonAnimationDefault() ButtonAnimation {
	return NewButtonAnimationScale(.98)
}

func NewButtonAnimationScale(v float32) ButtonAnimation {
	animationEnter := animation.NewAnimation(false,
		gween.NewSequence(
			gween.New(1, v, .1, ease.Linear),
		),
	)

	animationLeave := animation.NewAnimation(false,
		gween.NewSequence(
			gween.New(v, 1, .1, ease.Linear),
		),
	)

	animationClick := animation.NewAnimation(false,
		gween.NewSequence(
			gween.New(1, v, .1, ease.Linear),
			gween.New(v, 1, .4, ease.OutBounce),
		),
	)

	return ButtonAnimation{
		animationEnter: animationEnter,
		transformEnter: animation.TransformScaleCenter,
		animationLeave: animationLeave,
		transformLeave: animation.TransformScaleCenter,
		animationClick: animationClick,
		transformClick: animation.TransformScaleCenter,
	}
}

func NewButton(style ButtonStyle) *Button {
	return &Button{
		Style:         style,
		Clickable:     new(widget.Clickable),
		Label:         new(widget.Label),
		animClickable: new(widget.Clickable),

		Focused:          false,
		hoverSwitchState: false,
	}
}

func (btn *Button) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	animationEnter := btn.Style.Animation.animationEnter
	transformEnter := btn.Style.Animation.transformEnter
	animationLeave := btn.Style.Animation.animationLeave
	transformLeave := btn.Style.Animation.transformLeave
	animationClick := btn.Style.Animation.animationClick
	transformClick := btn.Style.Animation.transformClick

	clickable := btn.Clickable
	animClickable := btn.animClickable
	style := btn.Style

	return clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return animClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			semantic.Button.Add(gtx.Ops)

			{
				if animationEnter != nil {
					state := animationEnter.Update(gtx)
					if state.Active {
						transformEnter(gtx, state.Value).Add(gtx.Ops)
					}
				}
			}

			{
				if animationLeave != nil {
					state := animationLeave.Update(gtx)
					if state.Active {
						transformLeave(gtx, state.Value).Add(gtx.Ops)
					}
				}
			}

			{
				if animationClick != nil {
					state := animationClick.Update(gtx)
					if state.Active {
						transformClick(gtx, state.Value).Add(gtx.Ops)
					}
				}
			}

			backgroundColor := style.BackgroundColor
			textColor := style.TextColor

			if animClickable.Hovered() {
				pointer.CursorPointer.Add(gtx.Ops)
				if style.HoverBackgroundColor != nil {
					backgroundColor = *style.HoverBackgroundColor // f32color.Hovered(backgroundColor)
				}

				if style.HoverTextColor != nil {
					textColor = *style.HoverTextColor
				}
			}

			if animClickable.Hovered() && !btn.hoverSwitchState {
				btn.hoverSwitchState = true

				if animationEnter != nil {
					animationEnter.Start()
				}

				if animationLeave != nil {
					animationLeave.Reset()
				}
			}

			if !animClickable.Hovered() && btn.hoverSwitchState {
				btn.hoverSwitchState = false

				if animationLeave != nil {
					animationLeave.Start()
				}

				if animationEnter != nil {
					animationEnter.Reset()
				}
			}

			if animClickable.Clicked() {
				if animationClick != nil {
					animationClick.Reset().Start()
				}
			}

			c := op.Record(gtx.Ops)
			dims := style.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				if style.Icon != nil && btn.Style.Text == "" {
					return style.Icon.Layout(gtx, textColor)
				}

				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if style.Icon != nil {
							return style.Icon.Layout(gtx, textColor)
						}

						return layout.Dimensions{}
					}),
					layout.Rigid(layout.Spacer{Width: style.IconGap}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						paint.ColorOp{Color: textColor}.Add(gtx.Ops)
						return btn.Label.Layout(gtx, th.Shaper, style.Font,
							style.TextSize, style.Text, op.CallOp{})
					}),
				)
			})
			m := c.Stop()

			bounds := image.Rectangle{Max: dims.Size}
			paint.FillShape(gtx.Ops, backgroundColor,
				clip.UniformRRect(bounds, gtx.Dp(btn.Style.Rounded)).Op(gtx.Ops))

			m.Add(gtx.Ops)
			return dims
		})
	})
}
