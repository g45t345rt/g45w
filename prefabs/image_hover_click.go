package prefabs

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type ImageHoverClick struct {
	Image     *components.Image
	Clickable *widget.Clickable

	AnimationEnter   *animation.Animation
	AnimationLeave   *animation.Animation
	hoverSwitchState bool
}

func NewImageHoverClick(src image.Image) *ImageHoverClick {
	image := &components.Image{
		Src:     paint.NewImageOp(src),
		Fit:     components.Cover,
		Rounded: components.UniformRounded(unit.Dp(10)),
	}

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 1.1, .1, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1.1, 1, .1, ease.Linear),
	))

	return &ImageHoverClick{
		Image:          image,
		Clickable:      new(widget.Clickable),
		AnimationEnter: animationEnter,
		AnimationLeave: animationLeave,
	}
}

func (item *ImageHoverClick) Layout(gtx layout.Context) layout.Dimensions {
	return item.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		{
			state := item.AnimationEnter.Update(gtx)
			if state.Active {
				item.Image.Transform = func(dims layout.Dimensions, trans f32.Affine2D) f32.Affine2D {
					pt := dims.Size.Div(2)
					origin := f32.Pt(float32(pt.X), float32(pt.Y))
					return trans.Scale(origin, f32.Pt(state.Value, state.Value))
				}
			}
		}

		{
			state := item.AnimationLeave.Update(gtx)
			if state.Active {
				item.Image.Transform = func(dims layout.Dimensions, trans f32.Affine2D) f32.Affine2D {
					pt := dims.Size.Div(2)
					origin := f32.Pt(float32(pt.X), float32(pt.Y))
					return trans.Scale(origin, f32.Pt(state.Value, state.Value))
				}
			}
		}

		if item.Clickable.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
		}

		if item.Clickable.Hovered() && !item.hoverSwitchState {
			item.hoverSwitchState = true
			item.AnimationEnter.Start()
			item.AnimationLeave.Reset()
		}

		if !item.Clickable.Hovered() && item.hoverSwitchState {
			item.hoverSwitchState = false
			item.AnimationLeave.Start()
			item.AnimationEnter.Reset()
		}

		return item.Image.Layout(gtx)
	})
}
