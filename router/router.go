package router

import (
	"log"
	"sort"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type LayoutFunc func(gtx layout.Context, th *material.Theme)

type KeyLayout struct {
	DrawIndex int
	Layout    LayoutFunc
}

type Router struct {
	Pages   map[interface{}]Page // does not keep ordering with range (use drawOrder)
	Current interface{}

	drawOrder     []interface{}
	keyLayouts    []KeyLayout
	closeKeyboard bool
}

func NewRouter() *Router {
	return &Router{
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
		//if router.Current == tag {
		//return
		//}

		if router.Current != tag && router.Current != nil {
			router.Pages[router.Current].Leave()
		}

		router.closeKeyboard = true
		router.Current = tag
		router.Pages[router.Current].Enter()
	} else {
		log.Fatalf("container does not exists [%s]", tag)
	}
}

func (router *Router) AddLayout(keyLayout KeyLayout) {
	router.keyLayouts = append(router.keyLayouts, keyLayout)
	sort.Slice(router.keyLayouts, func(i, j int) bool {
		return router.keyLayouts[i].DrawIndex < router.keyLayouts[j].DrawIndex
	})
}

func (router *Router) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if router.closeKeyboard {
		// mobile: force close keyboard on page change
		// we probably don't need this since the keyboard automatically close when an input lose focus
		key.SoftKeyboardOp{Show: false}.Add(gtx.Ops)
		router.closeKeyboard = false
	}

	for _, tag := range router.drawOrder {
		page := router.Pages[tag]
		if page.IsActive() {
			page.Layout(gtx, th)
		}
	}

	for _, keyLayout := range router.keyLayouts {
		keyLayout.Layout(gtx, th)
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}
