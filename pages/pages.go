package pages

import (
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/g45t345rt/g45w/animation"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

const (
	PAGE_SETTINGS      = "page_settings"
	PAGE_NODE          = "page_node"
	PAGE_WALLET        = "page_wallet"
	PAGE_WALLET_SELECT = "page_wallet_select"
)

type PageSectionAnimation struct {
	animationEnter *animation.Animation
	animationLeave *animation.Animation
}

func NewPageSectionAnimation() *PageSectionAnimation {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .5, ease.OutExpo),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .5, ease.OutExpo),
	))

	return &PageSectionAnimation{
		animationEnter: animationEnter,
		animationLeave: animationLeave,
	}
}

func (p *PageSectionAnimation) Enter() bool {
	p.animationEnter.Start()
	p.animationLeave.Reset()
	return true
}

func (p *PageSectionAnimation) Leave() bool {
	p.animationEnter.Reset()
	p.animationLeave.Start()
	return true
}

func (p *PageSectionAnimation) Update(gtx layout.Context, finished func()) (trans op.TransformOp) {
	{
		state := p.animationEnter.Update(gtx)
		if state.Active {
			trans = animation.TransformY(gtx, state.Value)
		}
	}

	{
		state := p.animationLeave.Update(gtx)

		if state.Active {
			trans = animation.TransformY(gtx, state.Value)
		}

		if state.Finished {
			finished()
			op.InvalidateOp{}.Add(gtx.Ops)
		}
	}

	return
}
