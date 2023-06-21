package router

import (
	"log"

	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/olebedev/emitter"
)

type OnNavigate func(tag interface{})

var EVENT_CHANGE = "change"

type Router struct {
	Pages   map[interface{}]Page // does not keep ordering with range (use drawOrder)
	Current interface{}
	Event   *emitter.Emitter

	drawOrder []interface{}
	layouts   []func(gtx layout.Context, th *material.Theme)
}

func NewRouter() *Router {
	event := &emitter.Emitter{}
	event.Use("*", emitter.Void)
	return &Router{
		Event:     event,
		drawOrder: make([]interface{}, 0),
		Pages:     make(map[interface{}]Page),
	}
}

func (router *Router) Add(tag interface{}, page Page) {
	router.Pages[tag] = page
	router.drawOrder = append(router.drawOrder, tag)
}

func (router *Router) SetCurrent(tag interface{}) {
	_, ok := router.Pages[tag]
	if ok {
		if router.Current == tag {
			return
		}

		if router.Current != nil {
			router.Pages[router.Current].Leave()
		}

		router.Event.Emit(EVENT_CHANGE, tag)
		router.Current = tag
		router.Pages[router.Current].Enter()
	} else {
		log.Fatalf("container does not exists [%s]", tag)
	}
}

func (router *Router) PushLayout(layout func(gtx layout.Context, th *material.Theme)) {
	router.layouts = append(router.layouts, layout)
}

func (router *Router) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for _, tag := range router.drawOrder {
		page := router.Pages[tag]
		if page.IsActive() {
			page.Layout(gtx, th)
		}
	}

	for _, layout := range router.layouts {
		layout(gtx, th)
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}
