package prompt_modal

import (
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
)

type PromptModal struct {
	modal      *components.Modal
	input      *components.Input
	textHint   string
	submitChan chan string
}

var Instance *PromptModal

func LoadInstance() {
	modal := components.NewModal(components.ModalStyle{
		CloseOnOutsideClick: true,
		CloseOnInsideClick:  false,
		Direction:           layout.Center,
		Rounded: components.Rounded{
			SW: unit.Dp(10), SE: unit.Dp(10),
			NW: unit.Dp(10), NE: unit.Dp(10),
		},
		Animation: components.NewModalAnimationScaleBounce(),
	})

	input := components.NewInput()
	input.Border = widget.Border{}
	input.Inset = layout.Inset{}
	input.TextSize = unit.Sp(20)

	Instance = &PromptModal{
		modal: modal,
		input: input,
	}

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 3,
		Layout: func(gtx layout.Context, th *material.Theme) {
			Instance.Layout(gtx, th)
		},
	})
}

func (c *PromptModal) Open(text string, textHint string, inputHint key.InputHint) chan string {
	c.input.SetValue(text)
	c.input.Editor.SetCaret(len(text), len(text))
	c.input.Editor.InputHint = inputHint
	c.modal.SetVisible(true)
	c.input.Editor.Focus()
	c.textHint = textHint
	c.submitChan = make(chan string)
	return c.submitChan
}

func (c *PromptModal) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for _, e := range c.input.Editor.Events() {
		e, ok := e.(widget.SubmitEvent)
		if ok && e.Text != "" {
			c.input.SetValue("")
			c.submitChan <- e.Text
			c.modal.SetVisible(false)
			close(c.submitChan)
		}
	}

	c.modal.Style.Colors = theme.Current.ModalColors
	return c.modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(20), Bottom: unit.Dp(20),
			Left: unit.Dp(20), Right: unit.Dp(20),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Dp(200)
			gtx.Constraints.Max.X = gtx.Dp(200)
			c.input.Colors = theme.Current.InputColors
			return c.input.Layout(gtx, th, c.textHint)
		})
	})
}
