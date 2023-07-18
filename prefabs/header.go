package prefabs

import (
	"image/color"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/router"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Header struct {
	LabelTitle  material.LabelStyle
	Subtitle    layout.Widget
	ButtonRight *components.Button

	buttonGoBack *components.Button
	router       *router.Router
	history      []interface{}
}

func NewHeader(labelTitle material.LabelStyle, r *router.Router) *Header {
	textColor := color.NRGBA{R: 0, G: 0, B: 0, A: 100}
	textHoverColor := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

	walletIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	buttonGoBack := components.NewButton(components.ButtonStyle{
		Icon:           walletIcon,
		TextColor:      textColor,
		HoverTextColor: &textHoverColor,
	})

	header := &Header{
		LabelTitle:   labelTitle,
		buttonGoBack: buttonGoBack,
		router:       r,
		history:      make([]interface{}, 0),
	}

	return header
}

func (h *Header) SetTitle(title string) {
	h.LabelTitle.Text = title
}

func (h *Header) History() []interface{} {
	return h.history
}

func (h *Header) IsHistory(tag interface{}) bool {
	lastHistory := h.GetLastHistory()
	return tag == lastHistory
}

func (h *Header) AddHistory(tag interface{}) {
	if !h.IsHistory(tag) {
		h.history = append(h.history, tag)
	}
}

func (h *Header) GetLastHistory() interface{} {
	if len(h.history) == 0 {
		return nil
	}

	return h.history[len(h.history)-1]
}

func (h *Header) ResetHistory() {
	h.history = make([]interface{}, 0)
}

func (h *Header) GoBack() {
	if len(h.history) >= 2 {
		tag := h.history[len(h.history)-2]
		h.router.SetCurrent(tag)
		h.history = h.history[:len(h.history)-1]
	}
}

func (h *Header) handleKeyBack(gtx layout.Context) {
	key.InputOp{
		Tag:  h,
		Keys: key.NameEscape + "|" + key.NameBack,
	}.Add(gtx.Ops)

	for _, e := range gtx.Events(h) {
		switch e := e.(type) {
		case key.Event:
			if e.Name == key.NameEscape || e.Name == key.NameBack && e.State == key.Press { // don't use key.Release not implement on Android
				h.GoBack()
			}
		}
	}
}

func (h *Header) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	h.handleKeyBack(gtx)

	if h.buttonGoBack.Clickable.Clicked() {
		h.GoBack()
	}

	showBackButton := len(h.history) > 1

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if showBackButton {
				gtx.Constraints.Min.X = gtx.Dp(30)
				gtx.Constraints.Min.Y = gtx.Dp(30)
			} else {
				gtx.Constraints.Max.X = 0
				gtx.Constraints.Max.Y = 0
			}
			return h.buttonGoBack.Layout(gtx, th)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if showBackButton {
				return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
			}

			return layout.Dimensions{}
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return h.LabelTitle.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if h.Subtitle != nil {
						return h.Subtitle(gtx)
					}
					return layout.Dimensions{}
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if h.ButtonRight != nil {
				gtx.Constraints.Min.X = gtx.Dp(30)
				gtx.Constraints.Min.Y = gtx.Dp(30)
				return h.ButtonRight.Layout(gtx, th)
			}

			return layout.Dimensions{}
		}),
	)
}
