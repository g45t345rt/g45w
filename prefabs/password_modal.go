package prefabs

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/lang"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PasswordModal struct {
	Input              *components.Input
	animationWrongPass *animation.Animation
	iconLock           *widget.Icon

	Modal *components.Modal
}

func NewPasswordModal() *PasswordModal {
	input := components.NewPasswordInput()
	input.Border = widget.Border{}
	input.Inset = layout.Inset{}

	animationWrongPass := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .05, ease.Linear),
		gween.New(1, -1, .05, ease.Linear),
		gween.New(-1, 0, .05, ease.Linear),
	))

	iconLock, _ := widget.NewIcon(icons.ActionLock)

	modal := components.NewModal(components.ModalStyle{
		CloseOnOutsideClick: true,
		CloseOnInsideClick:  false,
		Direction:           layout.Center,
		BgColor:             color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Rounded:             components.UniformRounded(unit.Dp(10)),
		Inset:               layout.UniformInset(25),
		Animation:           components.NewModalAnimationScaleBounce(),
		Backdrop:            components.NewModalBackground(),
	})

	return &PasswordModal{
		Input:              input,
		Modal:              modal,
		animationWrongPass: animationWrongPass,
		iconLock:           iconLock,
	}
}

func (w *PasswordModal) StartWrongPassAnimation() {
	w.animationWrongPass.Reset()
	w.animationWrongPass.Start()
}

func (w *PasswordModal) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if w.Modal.Visible {
		if !w.Input.Editor.Focused() {
			w.Input.Editor.Focus()
		}
	} else {
		w.Input.SetValue("")
	}

	return w.Modal.Layout(gtx,
		func(gtx layout.Context) {
			{
				state := w.animationWrongPass.Update(gtx)
				if state.Active {
					value := float32(gtx.Dp(unit.Dp(state.Value * 50)))
					transform := f32.Affine2D{}.Offset(f32.Pt(value, 0))
					op.Affine(transform).Add(gtx.Ops)
				}
			}
		},
		func(gtx layout.Context) layout.Dimensions {
			// can't get capslock state with gioui - need to be implemented cross-platform
			/*if capsLockOn {
				offset := f32.Point{X: float32(gtx.Dp(30)), Y: float32(gtx.Dp(10))}
				trans := op.Affine(f32.Affine2D{}.Offset(offset)).Push(gtx.Ops)
				lbl := material.Label(th, unit.Sp(10), lang.Translate("CAPS LOCK IS ON"))
				lbl.Font.Weight = font.Bold
				lbl.Layout(gtx)
				trans.Pop()
			}*/

			return layout.UniformInset(unit.Dp(25)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.X = gtx.Dp(25)
						gtx.Constraints.Max.Y = gtx.Dp(25)
						return w.iconLock.Layout(gtx, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
						w.Input.TextSize = unit.Sp(20)
						return w.Input.Layout(gtx, th, lang.Translate("Enter password"))
					}),
				)
			})
		})
}
