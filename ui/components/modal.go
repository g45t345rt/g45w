package components

import (
	"image"
	"image/color"

	"gioui.org/io/pointer"
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
	var img *Image
	return func(gtx layout.Context) layout.Dimensions {
		if img == nil {
			src := utils.NewImageColor(gtx.Constraints.Max, color.RGBA{A: 100})
			img = &Image{
				Src: paint.NewImageOp(src),
			}
		}

		return img.Layout(gtx)
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
	Style        ModalStyle
	visible      bool
	clickableOut *widget.Clickable
	clickableIn  *widget.Clickable
}

func NewModal(style ModalStyle) *Modal {
	return &Modal{
		Style:        style,
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

		modal.Style.Animation.animationEnter.Start()
		modal.Style.Animation.animationLeave.Reset()
	} else {
		modal.Style.Animation.animationLeave.Start()
		modal.Style.Animation.animationEnter.Reset()
	}

	op.InvalidateOp{}.Add(gtx.Ops)
}

func (modal *Modal) Layout(gtx layout.Context, beforeDraw func(gtx layout.Context), w layout.Widget) layout.Dimensions {
	if !modal.visible {
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}

	animationEnter := modal.Style.Animation.animationEnter
	transformEnter := modal.Style.Animation.transformEnter
	animationLeave := modal.Style.Animation.animationLeave
	transformLeave := modal.Style.Animation.transformLeave

	clickedOut := modal.clickableOut.Clicked()
	clickedIn := modal.clickableIn.Clicked()

	if modal.Style.CloseOnOutsideClick {
		if clickedOut && !clickedIn {
			animationLeave.Start()
		}

		// I think its weird for outside click to see the pointer
		//if modal.clickableOut.Hovered() {
		//	pointer.CursorPointer.Add(gtx.Ops)
		//}
	}

	if modal.Style.CloseOnInsideClick {
		if clickedIn {
			animationLeave.Start()
		}

		if modal.clickableIn.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
		}
	}

	if modal.Style.Background != nil {
		modal.Style.Background(gtx)
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
		return modal.Style.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return modal.Style.Direction.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
