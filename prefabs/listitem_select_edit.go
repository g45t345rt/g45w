package prefabs

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type ListItemSelectEdit struct {
	ButtonSelect *components.Button
	ButtonEdit   *components.Button

	visible        bool
	animationEnter *animation.Animation
	animationLeave *animation.Animation
}

func NewListItemSelectEdit() *ListItemSelectEdit {
	buttonSelect := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonSelect.Label.Alignment = text.Middle
	buttonSelect.Style.Font.Weight = font.Bold

	buttonEdit := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonEdit.Label.Alignment = text.Middle
	buttonEdit.Style.Font.Weight = font.Bold

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .15, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .15, ease.Linear),
	))

	return &ListItemSelectEdit{
		ButtonEdit:     buttonEdit,
		ButtonSelect:   buttonSelect,
		animationEnter: animationEnter,
		animationLeave: animationLeave,
	}
}

func (n *ListItemSelectEdit) EditClicked() bool {
	return n.ButtonEdit.Clickable.Clicked()
}

func (n *ListItemSelectEdit) SelectClicked() bool {
	return n.ButtonSelect.Clickable.Clicked()
}

func (n *ListItemSelectEdit) Toggle() {
	n.SetVisible(!n.visible)
}

func (n *ListItemSelectEdit) SetVisible(visible bool) {
	if visible {
		n.visible = true
		n.animationEnter.Start()
		n.animationLeave.Reset()
	} else {
		n.animationEnter.Reset()
		n.animationLeave.Start()
	}
}

func (n *ListItemSelectEdit) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if !n.visible {
		return layout.Dimensions{}
	}

	{
		state := n.animationEnter.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	{
		state := n.animationLeave.Update(gtx)

		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}

		if state.Finished {
			n.visible = false
			op.InvalidateOp{}.Add(gtx.Ops)
		}
	}

	return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				n.ButtonSelect.Text = lang.Translate("SELECT")
				return n.ButtonSelect.Layout(gtx, th)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				n.ButtonEdit.Text = lang.Translate("EDIT")
				return n.ButtonEdit.Layout(gtx, th)
			}),
		)
	})
}
