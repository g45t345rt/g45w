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

func NewModalAnimationUp() ModalAnimation {
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

func NewModalAnimationDown() ModalAnimation {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.InCubic),
	))

	return ModalAnimation{
		animationEnter: animationEnter,
		transformEnter: animation.TransformY,
		animationLeave: animationLeave,
		transformLeave: animation.TransformY,
	}
}

func NewModalBackground() layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		src := utils.NewImageColor(gtx.Constraints.Max, color.RGBA{A: 100})
		image := Image{
			Src: paint.NewImageOp(src),
		}
		return image.Layout(gtx)
	}
}

type ModalStyle struct {
	CloseOnOutsideClick bool
	CloseOnInsideClick  bool
	Direction           layout.Direction
	Inset               layout.Inset
	Background          layout.Widget
	Animation           ModalAnimation
}

type Modal struct {
	ModalStyle   ModalStyle
	visible      bool
	clickableOut *widget.Clickable
	clickableIn  *widget.Clickable
}

func NewModal(style ModalStyle) *Modal {
	//	img := utils.NewImageColor(gtx.Constraints.Max, color.RGBA{A: 100})
	return &Modal{
		ModalStyle:   style,
		visible:      false,
		clickableOut: new(widget.Clickable),
		clickableIn:  new(widget.Clickable),
	}
}

func (modal *Modal) Visible() bool {
	return modal.visible
}

func (modal *Modal) SetVisible(gtx layout.Context, visible bool) {
	if visible {
		modal.visible = true

		modal.ModalStyle.Animation.animationEnter.Start()
		modal.ModalStyle.Animation.animationLeave.Reset()
	} else {
		modal.ModalStyle.Animation.animationLeave.Start()
		modal.ModalStyle.Animation.animationEnter.Reset()
	}

	op.InvalidateOp{}.Add(gtx.Ops)
}

func (modal *Modal) Layout(gtx layout.Context, beforeDraw func(gtx layout.Context), w layout.Widget) layout.Dimensions {
	if !modal.visible {
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}

	animationEnter := modal.ModalStyle.Animation.animationEnter
	transformEnter := modal.ModalStyle.Animation.transformEnter
	animationLeave := modal.ModalStyle.Animation.animationLeave
	transformLeave := modal.ModalStyle.Animation.transformLeave

	clickedOut := modal.clickableOut.Clicked()
	clickedIn := modal.clickableIn.Clicked()

	if modal.ModalStyle.CloseOnOutsideClick && clickedOut && !clickedIn {
		animationLeave.Start()
	}

	if modal.ModalStyle.CloseOnInsideClick && clickedIn {
		animationLeave.Start()
	}

	if modal.ModalStyle.Background != nil {
		modal.ModalStyle.Background(gtx)
	}

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
				op.InvalidateOp{}.Add(gtx.Ops)
				return layout.Dimensions{Size: gtx.Constraints.Max}
			}
		}
	}

	return modal.clickableOut.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return modal.ModalStyle.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return modal.ModalStyle.Direction.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				macro := op.Record(gtx.Ops)
				dims := modal.clickableIn.Layout(gtx, w)
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
