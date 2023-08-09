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
	"github.com/g45t345rt/g45w/animation"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type ButtonAnimation struct {
	animationEnter   *animation.Animation
	transformEnter   animation.TransformFunc
	animationLeave   *animation.Animation
	transformLeave   animation.TransformFunc
	animationClick   *animation.Animation
	transformClick   animation.TransformFunc
	animationLoading *animation.Animation
	transformLoading animation.TransformFunc
}

type ButtonColors struct {
	TextColor            color.NRGBA
	BackgroundColor      color.NRGBA
	HoverBackgroundColor *color.NRGBA
	HoverTextColor       *color.NRGBA
	BorderColor          color.NRGBA
}

type ButtonStyle struct {
	Rounded     Rounded
	TextSize    unit.Sp
	Inset       layout.Inset
	Font        font.Font
	Icon        *widget.Icon
	IconGap     unit.Dp
	Animation   ButtonAnimation
	Border      widget.Border
	LoadingIcon *widget.Icon
	Colors      ButtonColors
}

type Button struct {
	Text      string
	Style     ButtonStyle
	Clickable *widget.Clickable
	Label     *widget.Label
	Focused   bool
	Disabled  bool
	Loading   bool

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

	animationLoading := animation.NewAnimation(false,
		gween.NewSequence(
			gween.New(0, 1, 1, ease.Linear),
		),
	)
	animationLoading.Sequence.SetLoop(-1)

	return ButtonAnimation{
		animationEnter:   animationEnter,
		transformEnter:   animation.TransformScaleCenter,
		animationLeave:   animationLeave,
		transformLeave:   animation.TransformScaleCenter,
		animationClick:   animationClick,
		transformClick:   animation.TransformScaleCenter,
		animationLoading: animationLoading,
		transformLoading: animation.TransformRotate,
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

func (btn *Button) SetLoading(loading bool) {
	btn.Loading = loading
	btn.Disabled = loading

	animationLoading := btn.Style.Animation.animationLoading
	if loading {
		animationLoading.Reset().Start()
	} else {
		animationLoading.Pause()
	}
}

func (btn *Button) Clicked() bool {
	if btn.Disabled {
		return false
	}

	return btn.Clickable.Clicked()
}

func (btn *Button) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	animationEnter := btn.Style.Animation.animationEnter
	transformEnter := btn.Style.Animation.transformEnter
	animationLeave := btn.Style.Animation.animationLeave
	transformLeave := btn.Style.Animation.transformLeave
	animationClick := btn.Style.Animation.animationClick
	transformClick := btn.Style.Animation.transformClick
	animationLoading := btn.Style.Animation.animationLoading
	transformLoading := btn.Style.Animation.transformLoading

	clickable := btn.Clickable
	animClickable := btn.animClickable
	style := btn.Style
	colors := btn.Style.Colors

	return clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return animClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			semantic.Button.Add(gtx.Ops)

			{
				if animationEnter != nil {
					state := animationEnter.Update(gtx)
					if state.Active {
						defer transformEnter(gtx, state.Value).Push(gtx.Ops).Pop()
					}
				}
			}

			{
				if animationLeave != nil {
					state := animationLeave.Update(gtx)
					if state.Active {
						defer transformLeave(gtx, state.Value).Push(gtx.Ops).Pop()
					}
				}
			}

			{
				if animationClick != nil {
					state := animationClick.Update(gtx)
					if state.Active {
						defer transformClick(gtx, state.Value).Push(gtx.Ops).Pop()
					}
				}
			}

			backgroundColor := colors.BackgroundColor
			textColor := colors.TextColor

			if !btn.Disabled {
				if animClickable.Hovered() {
					pointer.CursorPointer.Add(gtx.Ops)
					if colors.HoverBackgroundColor != nil {
						backgroundColor = *colors.HoverBackgroundColor
					}

					if colors.HoverTextColor != nil {
						textColor = *colors.HoverTextColor
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
			} else {
				animationLeave.Reset()
				animationEnter.Reset()
				animationClick.Reset()
			}

			c := op.Record(gtx.Ops)
			style.Border.Color = colors.BorderColor
			dims := style.Border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return style.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					if style.Icon != nil && btn.Text == "" {
						return style.Icon.Layout(gtx, textColor)
					}

					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if style.LoadingIcon != nil && btn.Loading {
								r := op.Record(gtx.Ops)
								dims := style.LoadingIcon.Layout(gtx, textColor)
								c := r.Stop()

								{
									gtx.Constraints.Min = dims.Size

									if animationLoading != nil {
										state := animationLoading.Update(gtx)
										if state.Active {
											defer transformLoading(gtx, state.Value).Push(gtx.Ops).Pop()
										}
									}
								}

								c.Add(gtx.Ops)
								return dims
							} else if style.Icon != nil {
								return style.Icon.Layout(gtx, textColor)
							} else {
								return layout.Dimensions{}
							}
						}),
						layout.Rigid(layout.Spacer{Width: style.IconGap}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							paint.ColorOp{Color: textColor}.Add(gtx.Ops)
							return btn.Label.Layout(gtx, th.Shaper, style.Font,
								style.TextSize, btn.Text, op.CallOp{})
						}),
					)
				})
			})
			m := c.Stop()

			bounds := image.Rectangle{Max: dims.Size}
			paint.FillShape(gtx.Ops, backgroundColor,
				clip.RRect{
					Rect: bounds,
					SE:   gtx.Dp(style.Rounded.SE),
					SW:   gtx.Dp(style.Rounded.SW),
					NE:   gtx.Dp(style.Rounded.NE),
					NW:   gtx.Dp(style.Rounded.NW),
				}.Op(gtx.Ops),
			)

			m.Add(gtx.Ops)
			return dims
		})
	})
}
