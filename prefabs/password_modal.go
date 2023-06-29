package prefabs

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PasswordModal struct {
	editorStyle        material.EditorStyle
	animationWrongPass *animation.Animation
	iconLock           *widget.Icon

	Modal      *components.Modal
	submitted  bool
	submitText string
}

func NewPasswordModal() *PasswordModal {
	th := app_instance.Theme
	editor := new(widget.Editor)
	editor.SingleLine = true
	editor.Submit = true
	editor.Mask = rune(42)
	editorStyle := material.Editor(th, editor, lang.Translate("Enter password"))
	editorStyle.TextSize = unit.Sp(20)

	animationWrongPass := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .05, ease.Linear),
		gween.New(1, -1, .05, ease.Linear),
		gween.New(-1, 0, .05, ease.Linear),
	))

	iconLock, _ := widget.NewIcon(icons.ActionLock)

	w := app_instance.Window
	modal := components.NewModal(w, components.ModalStyle{
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
		editorStyle:        editorStyle,
		Modal:              modal,
		animationWrongPass: animationWrongPass,
		iconLock:           iconLock,
	}
}

func (w *PasswordModal) Submit() (bool, string) {
	if w.submitted {
		w.submitted = false
		return true, w.submitText
	}

	return false, w.submitText
}

func (w *PasswordModal) StartWrongPassAnimation() {
	w.animationWrongPass.Start()
}

func (w *PasswordModal) Layout(gtx layout.Context) layout.Dimensions {
	for _, e := range w.editorStyle.Editor.Events() {
		e, ok := e.(widget.SubmitEvent)
		if ok {
			//w.animationWrongPass.Start()
			w.editorStyle.Editor.SetText("")
			w.submitText = e.Text
			w.submitted = true
		}
	}

	if w.Modal.Visible() {
		w.editorStyle.Editor.Focus()
	}

	return w.Modal.Layout(gtx,
		func(gtx layout.Context) {
			{
				state := w.animationWrongPass.Update(gtx)
				if state.Active {
					transform := f32.Affine2D{}.Offset(f32.Pt(state.Value*15, 0))
					op.Affine(transform).Add(gtx.Ops)
				}
			}
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(25)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.X = gtx.Dp(25)
						gtx.Constraints.Max.Y = gtx.Dp(25)
						return w.iconLock.Layout(gtx, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
						return w.editorStyle.Layout(gtx)
					}),
				)
			})
		})
}
