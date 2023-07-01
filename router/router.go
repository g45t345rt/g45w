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

type SortKeyLayout []KeyLayout

func (k SortKeyLayout) Len() int           { return len(k) }
func (k SortKeyLayout) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }
func (k SortKeyLayout) Less(i, j int) bool { return k[i].DrawIndex < k[j].DrawIndex }

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
	sort.Sort(SortKeyLayout(router.keyLayouts))
}

func (router *Router) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if router.closeKeyboard {
		// mobile: force close keyboard on page change
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
