package components

import (
	"image"
	"time"

	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
)

type HoldPress struct {
	TriggerDuration time.Duration
	Triggered       bool

	hold      gesture.Click
	pressTime *time.Time
}

func NewHoldPress(duration time.Duration) *HoldPress {
	return &HoldPress{TriggerDuration: duration}
}

func (h *HoldPress) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	h.Triggered = false

	if h.pressTime != nil {
		if h.pressTime.Add(h.TriggerDuration).Before(gtx.Now) {
			h.pressTime = nil
			h.Triggered = true
		}
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	for _, e := range h.hold.Update(gtx.Queue) {
		switch e.Kind {
		case gesture.KindPress:
			h.pressTime = &gtx.Now
		case gesture.KindCancel, gesture.KindClick:
			h.pressTime = nil
		}
		op.InvalidateOp{}.Add(gtx.Ops) // make sure you invalidate or it won't trigger render sometimes (depends on mouse move if I have to guess)
	}

	r := op.Record(gtx.Ops)
	dims := w(gtx)
	c := r.Stop()

	defer clip.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops).Pop()
	h.hold.Add(gtx.Ops)
	c.Add(gtx.Ops)

	return dims
}
