package components

import (
	"image"
	"image/color"

	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type Rounded struct {
	NW unit.Dp
	NE unit.Dp
	SW unit.Dp
	SE unit.Dp
}

func UniformRounded(r unit.Dp) Rounded {
	return Rounded{
		NW: r, NE: r, SW: r, SE: r,
	}
}

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
		paint.ColorOp{Color: color.NRGBA{A: 100}}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}
}

type ModalStyle struct {
	CloseOnOutsideClick bool
	CloseOnInsideClick  bool
	Direction           layout.Direction
	Inset               layout.Inset
	Rounded             Rounded
	BgColor             color.NRGBA
	Backdrop            layout.Widget
	Animation           ModalAnimation
	CloseKeySet         key.Set
}

type Modal struct {
	Style        ModalStyle
	visible      bool
	clickableOut *widget.Clickable
	clickableIn  *widget.Clickable
}

func NewModal(style ModalStyle) *Modal {
	if style.CloseKeySet == "" {
		style.CloseKeySet = key.NameEscape + "|" + key.NameBack
	}

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

func (modal *Modal) handleKeyClose(gtx layout.Context) {
	key.InputOp{
		Tag:  modal,
		Keys: modal.Style.CloseKeySet,
	}.Add(gtx.Ops)

	for _, e := range gtx.Events(modal) {
		switch e := e.(type) {
		case key.Event:
			if e.State == key.Press {
				modal.SetVisible(gtx, false)
			}
		}
	}
}

func (modal *Modal) Layout(gtx layout.Context, beforeLayout func(gtx layout.Context), w layout.Widget) layout.Dimensions {
	if !modal.visible {
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}

	modal.handleKeyClose(gtx)

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

	if modal.Style.Backdrop != nil {
		modal.Style.Backdrop(gtx)
	}

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

			if state.Finished {
				modal.visible = false
				op.InvalidateOp{}.Add(gtx.Ops)
				//return layout.Dimensions{Size: gtx.Constraints.Max}
			}
		}
	}

	r := op.Record(gtx.Ops)
	dims := modal.Style.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return modal.Style.Direction.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			r := op.Record(gtx.Ops)
			dims := modal.clickableIn.Layout(gtx, w)
			c := r.Stop()

			if beforeLayout != nil {
				beforeLayout(gtx)
			}

			paint.FillShape(gtx.Ops, modal.Style.BgColor,
				clip.RRect{
					Rect: image.Rectangle{Max: dims.Size},
					SE:   gtx.Dp(modal.Style.Rounded.SE),
					SW:   gtx.Dp(modal.Style.Rounded.SW),
					NW:   gtx.Dp(modal.Style.Rounded.NW),
					NE:   gtx.Dp(modal.Style.Rounded.NE),
				}.Op(gtx.Ops),
			)

			c.Add(gtx.Ops)
			return dims
		})
	})
	c := r.Stop()

	if modal.Style.CloseOnOutsideClick {
		return modal.clickableOut.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			c.Add(gtx.Ops)
			return dims
		})
	}

	c.Add(gtx.Ops)
	return dims
}
