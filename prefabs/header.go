package prefabs

import (
	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Header struct {
	Title       func() string
	Subtitle    func(gtx layout.Context, th *material.Theme) layout.Dimensions
	ButtonRight *components.Button

	buttonGoBack *components.Button
	router       *router.Router
	history      []interface{}

	slideRight  gesture.Drag
	slideStartX float32
	calcSlide   bool
}

type HeaderTitleLayoutFunc func(gtx layout.Context, th *material.Theme, title string) layout.Dimensions

func NewHeader(r *router.Router) *Header {
	walletIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	buttonGoBack := components.NewButton(components.ButtonStyle{
		Icon: walletIcon,
	})

	header := &Header{
		buttonGoBack: buttonGoBack,
		router:       r,
		history:      make([]interface{}, 0),
	}

	return header
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

func (h *Header) handleKeyGoBack(gtx layout.Context) {
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

func (h *Header) swipeRightGoBack(gtx layout.Context) {
	var de *pointer.Event
	for _, e := range h.slideRight.Events(gtx.Metric, gtx, gesture.Horizontal) {
		switch e.Type {
		case pointer.Drag:
			de = &e
		case pointer.Press:
			h.slideStartX = 0
			h.calcSlide = true
		}
	}

	if de != nil && h.calcSlide {
		if h.slideStartX == 0 {
			h.slideStartX = de.Position.X
		}

		slideWidth := float32(gtx.Constraints.Max.Div(2).X) // half the screen
		if h.slideStartX-de.Position.X < -slideWidth {
			h.calcSlide = false
			h.GoBack()
		}
	}

	h.slideRight.Add(gtx.Ops)
}

func (h *Header) Layout(gtx layout.Context, th *material.Theme, titleLayout HeaderTitleLayoutFunc) layout.Dimensions {
	h.handleKeyGoBack(gtx)
	h.swipeRightGoBack(gtx)

	if h.buttonGoBack.Clicked() {
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
			h.buttonGoBack.Style.Colors = theme.Current.HeaderBackButtonColors
			return h.buttonGoBack.Layout(gtx, th)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if showBackButton {
				return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
			}

			return layout.Dimensions{}
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			var childs []layout.FlexChild

			childs = append(childs, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return titleLayout(gtx, th, h.Title())
			}))

			if h.Subtitle != nil {
				childs = append(childs, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return h.Subtitle(gtx, th)
				}))
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, childs...)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if h.ButtonRight != nil {
				gtx.Constraints.Min.X = gtx.Dp(30)
				gtx.Constraints.Min.Y = gtx.Dp(30)
				h.ButtonRight.Style.Colors = theme.Current.ButtonIconPrimaryColors
				return h.ButtonRight.Layout(gtx, th)
			}

			return layout.Dimensions{}
		}),
	)
}
