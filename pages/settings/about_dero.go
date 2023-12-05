package page_settings

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/lang"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type PageDero struct {
	isActive       bool
	list           *widget.List
	animationEnter *animation.Animation
	animationLeave *animation.Animation
	infoItems      []*InfoListItem
}

func NewPageDero() *PageDero {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	// do not remove @lang.Translate comment
	// it's used by the python script to generate language json dictionary
	// we don't use lang.Translate directly here because it needs to be inside the Layout func or the value won't be updated after language change
	infoItems := []*InfoListItem{
		NewInfoListItem("Website", "https://dero.io", text.WrapGraphemes),                      //@lang.Translate("App Directory")
		NewInfoListItem("Github", "https://github.com/deroproject/derohe", text.WrapGraphemes), //@lang.Translate("Wallets Directory")
		NewInfoListItem("Forum", "https://forum.dero.io", text.WrapGraphemes),                  //@lang.Translate("Cache Directory")
		NewInfoListItem("Docs", "https://docs.dero.io", text.WrapGraphemes),                    //@lang.Translate("Version")
	}

	return &PageDero{
		infoItems:      infoItems,
		list:           list,
		animationEnter: animationEnter,
		animationLeave: animationLeave,
	}
}

func (p *PageDero) IsActive() bool {
	return p.isActive
}

func (p *PageDero) Enter() {
	p.isActive = true
	page_instance.header.Title = func() string { return lang.Translate("About DERO") }

	if !page_instance.header.IsHistory(PAGE_APP_INFO) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageDero) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *PageDero) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	{
		state := p.animationEnter.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	{
		state := p.animationLeave.Update(gtx)
		if state.Finished {
			p.isActive = false
			op.InvalidateOp{}.Add(gtx.Ops)
		}

		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	var widgets []layout.Widget

	for i := range p.infoItems {
		idx := i
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return p.infoItems[idx].Layout(gtx, th)
		})
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
	})

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay
	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return widgets[index](gtx)
	})
}
