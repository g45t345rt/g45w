package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/utils"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type ModalAnimation struct {
	animationEnter *animation.Animation
	transformEnter animation.TransformFunc
	animationLeave *animation.Animation
	transformLeave animation.TransformFunc
}

func NewModalAnimationScaleBounce() ModalAnimation {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.OutBounce),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.OutBounce),
	))

	return ModalAnimation{
		animationEnter: animationEnter,
		transformEnter: animation.TransformScaleCenter,
		animationLeave: animationLeave,
		transformLeave: animation.TransformScaleCenter,
	}
}

func NewModalAnimationDownUp() ModalAnimation {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.InCubic),
	))

	return ModalAnimation{
		animationEnter: animationEnter,
		transformEnter: animation.TransformY,
		animationLeave: animationLeave,
		transformLeave: animation.TransformY,
	}
}

type Modal struct {
	visible      bool
	clickableOut *widget.Clickable
	clickableIn  *widget.Clickable
	direction    layout.Direction
	inset        layout.Inset
	bg           *Image

	animation ModalAnimation
}

func NewModal(th *material.Theme, direction layout.Direction, inset layout.Inset, modalAnimation ModalAnimation) *Modal {
	return &Modal{
		direction:    direction,
		inset:        inset,
		clickableOut: new(widget.Clickable),
		clickableIn:  new(widget.Clickable),
		animation:    modalAnimation,
	}
}

func (modal *Modal) Visible() bool {
	return modal.visible
}

func (modal *Modal) SetVisible(gtx layout.Context, visible bool) {
	if visible {
		modal.visible = true

		modal.animation.animationEnter.Start()
		modal.animation.animationLeave.Reset()
	} else {
		modal.animation.animationLeave.Start()
		modal.animation.animationEnter.Reset()
	}

	op.InvalidateOp{}.Add(gtx.Ops)
}

func (modal *Modal) Layout(gtx layout.Context, beforeDraw func(gtx layout.Context), widget layout.Widget) layout.Dimensions {
	if !modal.visible {
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}

	animationEnter := modal.animation.animationEnter
	transformEnter := modal.animation.transformEnter
	animationLeave := modal.animation.animationLeave
	transformLeave := modal.animation.transformLeave

	if modal.clickableOut.Clicked() && !modal.clickableIn.Clicked() {
		animationLeave.Start()
	}

	if modal.bg == nil {
		img := utils.NewImageColor(gtx.Constraints.Max, color.RGBA{R: 0, G: 0, B: 0, A: 100})
		modal.bg = &Image{
			Src: paint.NewImageOp(img),
		}
	}

	modal.bg.Layout(gtx)

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

			if state.Finished {
				modal.visible = false
				modal.bg = nil
				op.InvalidateOp{}.Add(gtx.Ops)
				return layout.Dimensions{Size: gtx.Constraints.Max}
			}
		}
	}

	return modal.clickableOut.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return modal.inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return modal.direction.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				macro := op.Record(gtx.Ops)
				dims := modal.clickableIn.Layout(gtx, widget)
				c := macro.Stop()

				if beforeDraw != nil {
					beforeDraw(gtx)
				}

				paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255},
					clip.UniformRRect(
						image.Rectangle{Max: dims.Size},
						gtx.Dp(unit.Dp(10)),
					).Op(gtx.Ops),
				)

				c.Add(gtx.Ops)
				return dims
			})
		})
	})
}
