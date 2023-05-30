package router

import (
	"log"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

type Router struct {
	Containers map[interface{}]Container // does not keep ordering with range (use drawOrder)
	Current    interface{}
	Primary    interface{}
	drawOrder  []interface{}

	layouts []func(gtx layout.Context, th *material.Theme)
}

func NewRouter() *Router {
	return &Router{
		Containers: make(map[interface{}]Container),
	}
}

func (router *Router) IsPrimary() bool {
	return router.Primary == router.Current
}

func (router *Router) Add(tag interface{}, container Container) {
	router.Containers[tag] = container
	router.drawOrder = append(router.drawOrder, tag)
}

func (router *Router) SetPrimary(tag interface{}) {
	router.Primary = tag
	router.SetCurrent(tag)
}

func (router *Router) SetCurrent(tag interface{}) {
	_, ok := router.Containers[tag]
	if ok {
		if router.Current == tag {
			return
		}

		if router.Current != nil {
			router.Containers[router.Current].Leave()
		}

		router.Current = tag
		router.Containers[router.Current].Enter()
	} else {
		log.Fatalf("container does not exists [%s]", tag)
	}
}

func (router *Router) PushLayout(layout func(gtx layout.Context, th *material.Theme)) {
	router.layouts = append(router.layouts, layout)
}

func (router *Router) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for _, tag := range router.drawOrder {
		page := router.Containers[tag]
		if page.IsActive() {
			page.Layout(gtx, th)
		}
	}

	for _, layout := range router.layouts {
		layout(gtx, th)
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}
