package prefabs

import (
	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Header struct {
	Title       func() string
	Subtitle    func(gtx layout.Context, th *material.Theme) layout.Dimensions
	LeftLayout  func(gtx layout.Context, th *material.Theme) layout.Dimensions
	RightLayout func(gtx layout.Context, th *material.Theme) layout.Dimensions

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

func (h *Header) GetPosition(tag interface{}) int {
	for i, history := range h.history {
		if history == tag {
			return i
		}
	}

	return -1
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

func (h *Header) HandleKeyGoBack(gtx layout.Context) {
	for _, e := range gtx.Events(h) {
		switch e := e.(type) {
		case key.Event:
			if e.Name == key.NameEscape || e.Name == key.NameBack && e.State == key.Press { // don't use key.Release not implement on Android
				h.GoBack()
			}
		}
	}

	key.InputOp{
		Tag:  h,
		Keys: key.NameEscape + "|" + key.NameBack,
	}.Add(gtx.Ops)
}

func (h *Header) HandleSwipeRightGoBack(gtx layout.Context) {
	// use swipe action only for mobile
	if !utils.IsMobile() {
		return
	}

	var de *pointer.Event
	for _, e := range h.slideRight.Update(gtx.Metric, gtx, gesture.Horizontal) {
		switch e.Kind {
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
	if h.buttonGoBack.Clicked(gtx) {
		h.GoBack()
	}

	showBackButton := len(h.history) > 1

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if showBackButton {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(30)
						gtx.Constraints.Min.Y = gtx.Dp(30)
						h.buttonGoBack.Style.Colors = theme.Current.HeaderBackButtonColors
						return h.buttonGoBack.Layout(gtx, th)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				)
			}

			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if h.LeftLayout != nil {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return h.LeftLayout(gtx, th)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				)
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

			return layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			}.Layout(gtx, childs...)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if h.RightLayout != nil {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return h.RightLayout(gtx, th)
					}),
				)
			}

			return layout.Dimensions{}
		}),
	)
}

/*
func GetPageHeaderAnimation(pos int) (enter *animation.Animation, leave *animation.Animation) {
	if pos%2 == 0 {
		enter = animation.NewAnimation(false, gween.NewSequence(
			gween.New(-1, 0, .25, ease.Linear),
		))
		leave = animation.NewAnimation(false, gween.NewSequence(
			gween.New(0, -1, .25, ease.Linear),
		))
	} else {
		enter = animation.NewAnimation(false, gween.NewSequence(
			gween.New(1, 0, .25, ease.Linear),
		))
		leave = animation.NewAnimation(false, gween.NewSequence(
			gween.New(0, 1, .25, ease.Linear),
		))
	}

	return
}*/

type PageHeaderAnimation struct {
	tag            interface{}
	animationEnter *animation.Animation
	animationLeave *animation.Animation
}

func NewPageHeaderAnimation(tag interface{}) *PageHeaderAnimation {
	return &PageHeaderAnimation{
		tag: tag,
	}
}

// This is a helper to load page animation based on the header depth
func (p *PageHeaderAnimation) Enter(header *Header) bool {
	pos := header.GetPosition(p.tag)
	if pos == -1 {
		pos = len(header.history)
	}

	if pos%2 == 0 {
		p.animationEnter = animation.NewAnimation(false, gween.NewSequence(
			gween.New(-1, 0, .25, ease.OutExpo),
		))
		p.animationLeave = animation.NewAnimation(false, gween.NewSequence(
			gween.New(0, -1, .25, ease.OutExpo),
		))
	} else {
		p.animationEnter = animation.NewAnimation(false, gween.NewSequence(
			gween.New(1, 0, .25, ease.OutExpo),
		))
		p.animationLeave = animation.NewAnimation(false, gween.NewSequence(
			gween.New(0, 1, .25, ease.OutExpo),
		))
	}

	// this is when a section enter and trigger the page enter event
	// we don't want to execute the animation
	if !header.IsHistory(p.tag) {
		p.animationLeave.Reset()
		p.animationEnter.Start()
	}

	return true
}

func (p *PageHeaderAnimation) Leave(header *Header) bool {
	// same here when a section leave we don't want to trigger the page animation
	if header.IsHistory(p.tag) {
		p.animationEnter.Reset()
		p.animationLeave.Start()
		return true
	} else {
		return false
		//a.page.SetActive(false)
	}
}

func (p *PageHeaderAnimation) Update(gtx layout.Context, finished func()) (trans op.TransformOp) {
	{
		// check if nil because it's not instanciated in constructor (just to be sure and avoid race condition with Enter func and Update loop)
		if p.animationEnter != nil {
			state := p.animationEnter.Update(gtx)
			if state.Active {
				trans = animation.TransformX(gtx, state.Value)
			}
		}
	}

	{
		if p.animationLeave != nil {
			state := p.animationLeave.Update(gtx)
			if state.Active {
				trans = animation.TransformX(gtx, state.Value)
			}

			if state.Finished {
				finished()
				//a.page.SetActive(false)
				op.InvalidateOp{}.Add(gtx.Ops)
			}
		}
	}

	return
}
